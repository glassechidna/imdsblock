# `imdsblock`

You might not want your Docker containers to have access to the host's EC2 
instance profile's IAM credentials. That said, you might still want it to have 
access to the rest of the EC2 instance metadata service. There's all sorts
of useful, mostly innocuous data in there, like the instance's availability zone.

`imdsblock` tries to achieve a compromise wherein requests for IAM credentials
from Docker containers are blocked, but everything else is permitted. This is
achieved through a combination of `iptables` rules and a reverse proxy listening
on `127.0.0.1:51999`.

## Easy installation

You can install an RPM compatible with Amazon Linux 2 as follows:

    yum install -y <url>

It will install the reverse proxy and a systemd unit file that runs the daemon 
and configures iptables rules on your behalf. Only installation is necessary, 
the post-install script will automatically start the `imdsblock` service.

## Manul

Either compile the binary yourself or download the tarball. You can then run it
however you choose and use the following `iptables` incantation:

    iptables -t nat -A PREROUTING -p tcp -d 169.254.169.254 --dport 80 -j DNAT --to-destination 127.0.0.1:51999
