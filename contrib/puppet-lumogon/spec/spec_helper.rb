require 'puppetlabs_spec_helper/module_spec_helper'
require 'simplecov'
require 'simplecov-console'

SimpleCov.start do
  add_filter '/spec'
  formatter SimpleCov::Formatter::MultiFormatter.new([
    SimpleCov::Formatter::HTMLFormatter,
    SimpleCov::Formatter::Console
  ])
end

RSpec.configure do |config|
  config.mock_with :rspec
end
