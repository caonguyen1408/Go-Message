What we have
------------
- RabbitMQ Server
- Redis Server
- Web Server
- BE Server

Required:
------
- Docker
- Ansible

Usage
------

#### Build docker
`cd ansible`
`ansible-playbook -e "ansible_ip=127.0.0.1" deploy.yml`
   + if you run localhost ansible_ip = 127.0.0.1
   + if you run VPS ansibe_ip = VPS_IP
- Get a cup of coffe and relax, deploy will take 5-10min


#### Test
`http://127.0.0.1/send.html`

