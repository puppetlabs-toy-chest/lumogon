from fixtures import lumogon

image_under_test = "debian"

def test_os(lumogon):
    family = lumogon['capabilities']['host']['payload']['platformfamily']
    assert "debian" == family

def test_systemd_present(lumogon):
    print lumogon
    packages = lumogon['capabilities']['dpkg']['payload']
    assert "systemd" in packages
