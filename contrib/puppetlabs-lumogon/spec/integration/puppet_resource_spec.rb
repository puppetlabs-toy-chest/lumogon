require 'spec_helper'

describe "puppetlabs-lumogon integration" do
  before(:all) do
    `docker run -it --privileged --name dind-test -d docker:dind`
    `docker run -it --rm --link dind-test:docker docker run -d ubuntu /bin/sh -c "while true; do echo hello world; sleep 1; done"`
  end

  after(:all) do
    `docker kill dind-test | xargs docker rm`
  end

  describe 'puppet resource lumogon' do
    command 'bundle exec puppet resource lumogon'
    environment RUBYLIB: $LOAD_PATH.join(':')
    its(:exitstatus) { is_expected.to eq 0 }
    its(:stdout) { is_expected.to include "lumogon { '/dind-test':" }
    its(:stdout) { is_expected.to include 'apk' }
    its(:stdout) { is_expected.to include 'os' }
    its(:stderr) { is_expected.to eq '' }
  end
end
