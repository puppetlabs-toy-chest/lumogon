require 'json'

Puppet::Type.type(:lumogon).provide(:lumogon) do
  desc "Provider for Lumogon"
  commands :docker => 'docker'

  mk_resource_methods

  def self.instances #rubocop:disable Metrics/AbcSize
    scan_report = docker(:run, '--rm', '-v', '/var/run/docker.sock:/var/run/docker.sock', 'puppet/lumogon', :scan)
    containers = JSON.parse(scan_report)['containers']

    containers.collect do |c|
      container = c[1]

      new({
        :name            => container['container_name'],

        # Host Capability
        :id              => container['container_id'],
        :hostname        => container['capabilities']['host']['payload']['hostname'],
        :os              => container['capabilities']['host']['payload']['os'],
        :platform        => container['capabilities']['host']['payload']['platform'],
        :platformfamily  => container['capabilities']['host']['payload']['platformfamily'],
        :platformversion => container['capabilities']['host']['payload']['platformversion'],

        # Package Capability
        :apk             => container['capabilities'].key?('apk') ? container['capabilities']['apk']['payload'] : :absent,
        :dpkg            => container['capabilities'].key?('dpkg') ? container['capabilities']['dpkg']['payload'] : :absent,
        :yum             => container['capabilities'].key?('yum') ? container['capabilities']['yum']['payload'] : :absent,

        # Label Capability
        :labels          => container['capabilities']['label'].key?('payload') ? container['capabilities']['label']['payload'] : :absent
      })
    end
  end
end
