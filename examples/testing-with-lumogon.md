# Testing with Lumogon

Given the ability to generate a data structure from a container or
image, a desirable use case for Lumogon is introducing tests to the image
build process, or building a constant-compliance style test suite for
production monitoring purposes. The following is a demonstration of how
simple this is building on top of Lumogon, more than it being a full
example. Depending on user feedback, and the interest in the open source
community, this would be a great area to explore and expand on.

Using less than 30 lines of Python we can create a framework to
write tests like the following. Because of how Lumogon uses data as a
first class interface it’s simple to imagine writing similar tests
in other languages or testing frameworks too.

```python
from fixtures import lumogon

image_under_test = "debian"

def test_os(lumogon):
    family = lumogon['capabilities']['host']['payload']['platformfamily']
    assert "debian" == family

def test_systemd_present(lumogon):
    packages = lumogon['capabilities']['dpkg']['payload']
    assert "systemd" in packages

def test_label_present(lumogon):
    labels = lumogon['capabilities']['label']['payload']
    assert "com.docker.compose.service" in labels
```

You'll find the full working example in the
[Lumogon repository](https://github.com/puppetlabs/lumogon/tree/master/examples).

The accompanying `Makefile` in the `examples/testing-with-lumogon`
directory should let you try this out. This requires Python and virtualenv
but otherwise you should just need to run:

```
make
```

Do try and write you're own assertions against the returned data, or
change the `image_under_test` to your own Docker image.

These tests could be used in a CI pipeline to verify that metadata was
set correctly, or that certain packages were installed, or that
known-vulnerable packages (across multiple different package managers or
operating systems) were not present, or that all images where using the
corporate standard operating system.

If you're interested in these examples or have any questions about
Lumogon then head over to our [Slack channel](https://puppetcommunity.slack.com/messages/G58F97FC5).
