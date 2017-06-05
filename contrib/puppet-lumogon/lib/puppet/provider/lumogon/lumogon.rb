require 'json'
require 'pry'

Puppet::Type.type(:lumogon).provide(:lumogon) do
  desc "Provider for Lumogon"

  commands :docker => '/usr/bin/docker'

  def self.instances
    scan_report = docker(:run, '--rm', '-v', '/var/run/docker.sock:/var/run/docker.sock', 'puppet/lumogon', :scan)
    containers = JSON.parse(scan_report)['containers']

    containers.each do |c|
      container = c[1]

      new({
        :container_name  => container['container_name'],
        :ensure          => :present,

        # Host Capability
        :id              => container['container_id'],
        :hostname        => container['capabilities']['host']['payload']['hostname'],
        :os              => container['capabilities']['host']['payload']['os'],
        :platform        => container['capabilities']['host']['payload']['platform'],
        :platformfamily  => container['capabilities']['host']['payload']['platformfamily'],
        :platformversion => container['capabilities']['host']['payload']['platformversion'],

        # Package Capability
        :apk             => container['capabilities'].has_key?('apk') ? container['capabilities']['apk']['payload'] : :absent,
        :dpkg            => container['capabilities'].has_key?('dpkg') ? container['capabilities']['dpkg']['payload'] : :absent,
        :yum             => container['capabilities'].has_key?('yum') ? container['capabilities']['yum']['payload'] : :absent,

        # Label Capability
        :labels          => container['capabilities']['label'].has_key?('payload') ? container['capabilities']['label']['payload'] : :absent
      })
    end
  end

  def name
    @resource[:name]
  end


end
