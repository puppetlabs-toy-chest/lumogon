require 'spec_helper_acceptance'

lumogon_command = [
  'docker',
  'run',
  '--rm',
  '-v',
  '/var/run/docker.sock:/var/run/docker.sock',
  'puppet/lumogon',
  '--disable-analytics'
].join(' ')

describe command(lumogon_command) do
  its(:exit_status) { should eq 0 }
end
