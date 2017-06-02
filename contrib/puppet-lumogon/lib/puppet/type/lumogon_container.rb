Puppet::Type.newtype(:lumogon_container) do
  newparam(:name) do
    isnamevar
  end

  # Host Capability
  newproperty(:id)
  newproperty(:hostname)
  newproperty(:os)
  newproperty(:platform)
  newproperty(:platformfamily)
  newproperty(:platformversion)

  # Package Capability
  newproperty(:apk)
  newproperty(:dpkg)
  newproperty(:yum)

  # Label Capability
  newproperty(:labels)
end
