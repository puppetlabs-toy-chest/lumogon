require 'beaker-rspec/spec_helper'
require 'beaker-rspec/helpers/serverspec'
require 'beaker/puppet_install_helper'
require 'beaker_spec_helper'

include BeakerSpecHelper

# automatically load any shared examples or contexts
Dir['./spec/support/**/*.rb'].sort.each { |f| require f }

ENV['PUPPET_INSTALL_TYPE'] = ENV['PUPPET_INSTALL_TYPE'] || 'agent'
run_puppet_install_helper unless ENV['BEAKER_provision'] == 'no'

RSpec.configure do |c|
  proj_root = File.expand_path(File.join(File.dirname(__FILE__), '..'))
  module_name = proj_root.split('-').last
  c.formatter = :documentation
  c.before :suite do
    puppet_module_install(source: proj_root, module_name: module_name)
    hosts.each do |host|
      BeakerSpecHelper::spec_prep(host)
    end
  end
end
