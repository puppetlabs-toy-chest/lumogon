require 'json'
require 'pry'

Puppet::Type.type(:lumogon).provide(:lumogon) do #rubocop:disable Metrics/BlockLength
  desc "Provider for Lumogon"
  commands :docker => 'docker'

  mk_resource_methods

  def self.instances #rubocop:disable Metrics/AbcSize
    cli_output = docker(:run, '--rm', '-v', '/var/run/docker.sock:/var/run/docker.sock', 'puppet/lumogon', :scan)
    scan_report = self.parse_lumogon_report(cli_output)
    containers = scan_report['containers']

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

  # Custom function to parse output from Lumogon as Puppet muxes together both
  # STDOUT and STDERR as part of Command Execution
  # https://github.com/puppetlabs/puppet/blob/master/lib/puppet/provider.rb#L259
  def self.parse_lumogon_report(report)
    # Collapse output string from Puppet (STDERR and STDOUT) into a single string
    # removing all whitespace from Lumogon Pretty Print report
    collapsed_output = report.split("\n").map(&:strip).join("")

    # Scan Output grab the STDOUT JSON Blob
    json_report = collapsed_output.match(/{.+}/)[0]

    # Parse and return
    JSON.parse(json_report)
  end
end
