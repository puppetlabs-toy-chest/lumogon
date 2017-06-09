Puppet::Type.newtype(:lumogon) do #rubocop:disable Metrics/BlockLength
  @doc = %q{Inspect running container with Lumogon}

  newparam(:name) do
    desc "The name of the container"
    isnamevar
  end

  # Host Capability
  newproperty(:id) do
    desc "The container ID of the running container"
  end

  newproperty(:hostname) do
    desc "The hostname of the running container"
  end

  newproperty(:os) do
   desc "The Operating system of the running container"
  end

  newproperty(:platform) do
    desc "The distribution platform of the running container"
  end

  newproperty(:platformfamily) do
    desc "The distribution family of the running container"
  end

  newproperty(:platformversion) do
    desc "The distribution version of the running container"
  end

  # Package Capability
  newproperty(:apk) do
    desc "A list of all Alpine Linu packages installed inside the running container"
  end

  newproperty(:dpkg) do
    desc "A list of all Debian packages installed inside the running container"
  end

  newproperty(:yum) do
    desc "A list of all RedHat packages installed inside the running container"
  end

  # Label Capability
  newproperty(:labels) do
    desc "A list of labels attached to the running container"
  end
end
