import pytest
import docker
import json

@pytest.fixture(scope="module")
def lumogon(request):
    client = docker.from_env()
    cmd = '/bin/sh -c "while true; do echo hello world; sleep 1; done"'
    image = getattr(request.module, "image_under_test", "ubuntu")
    temporary_container = client.containers.run(image, cmd, entrypoint="", detach=True)
    output = client.containers.run("puppet/lumogon",
                                   "scan %s" % temporary_container.id,
                                   volumes={'/var/run/docker.sock': {'bind': '/var/run/docker.sock', 'mode': 'rw'}},
                                   remove=True)
    data = json.loads(output)
    container_id = data['containers'].keys()[0]
    yield data['containers'][container_id]
    temporary_container.remove(force=True)
