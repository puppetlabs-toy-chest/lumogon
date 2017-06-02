require 'json'
require 'pry'

Puppet::Type.type(:lumogon_container).provide(:lumogon) do
  commands :docker => '/usr/bin/docker'
  mk_resource_methods

  def self.instances
    report = docker('run', '--rm', '-v', '/var/run/docker.sock:/var/run/docker.sock', 'puppet/lumogon', 'scan')
    containers = JSON.parse(report)['containers']

    containers.each do |container|
      results = container[1]
      new(
        :name            => results['container_name'],
        :id              => results['container_id'],
        :hostname        => results['capabilities']['host']['payload']['hostname'],
        :os              => results['capabilities']['host']['payload']['os'],
        :platform        => results['capabilities']['host']['payload']['platform'],
        :platformfamily  => results['capabilities']['host']['payload']['platformfamily'],
        :platformversion => results['capabilities']['host']['payload']['platformversion'],
        :apk             => results['capabilities'].has_key?('apk')  ? container['capabilities']['apk']['payload'] : :absent,
        :dpkg            => results['capabilities'].has_key?('dpkg') ? container['capabilities']['dpkg']['payload'] : :absent,
        :yum             => results['capabilities'].has_key?('yum')  ? container['capabilities']['yum']['payload'] : :absent,
        :labels          => results['capabilities']['labels'].has_key?('payload') ? container['capabilities']['labels']['payload'] : :absent
      )
    end
  end

  def exists?
  end

end
