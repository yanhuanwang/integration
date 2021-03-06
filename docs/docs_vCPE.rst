.. This work is licensed under a Creative Commons Attribution 4.0
   International License. http://creativecommons.org/licenses/by/4.0
   Copyright 2018 Huawei Technologies Co., Ltd.  All rights reserved.

.. _docs_vcpe:

vCPE Use Case
----------------------------

Description
~~~~~~~~~~~
vCPE use case is based on Network Enhanced Residential Gateway architecture specified in Technical Report 317 (TR-317), which defines how service providers deploy residential broadband services like High Speed Internet Access. The use case implementation has infrastructure services and customer service. The common infrastructure services are deployed first and shared by all customers. The use case demonstrates ONAP capabilities to design, deploy, configure and control sophisticated services.      

More details on the vCPE Use Case can be found on wiki page https://wiki.onap.org/pages/viewpage.action?pageId=3246168

Source Code
~~~~~~~~~~~
vcpe test scripts: https://git.onap.org/integration/tree/test/vcpe?h=dublin

How to Use
~~~~~~~~~~
Most part of the use case has been automated by vcpe scripts. For the details on how to run the scripts, please refer to the use case tutorial on https://wiki.onap.org/display/DW/vCPE+Use+Case+Tutorial%3A+Design+and+Deploy+based+on+ONAP.

Here are the main steps to run the use case in Integration lab environment, where vCPE script is pre-installed on Rancher node under /root/integration/test/vcpe:

1. Run Robot script from Rancher node to onboard VNFs, create and distribute models for vCPE four infrastructure services, i.e. infrastructure, brg, bng and gmux

:: 

   demo-k8s.sh onap init
 
2. Add customer SDN-ETHERNET-INTERNET (see the use case tutorial wiki page for detail)

3. Add identity-url to RegionOne data in A&AI. First use POSTMAN to GET cloud-region RegionOne data, then add identity-url and PUT back to A&AI

::

   GET https://{{aai}}:{{port}}/aai/v14/cloud-infrastructure/cloud-regions/cloud-region/CloudOwner/RegionOne

::

   PUT https://{{aai}}:{{port}}/aai/v14/cloud-infrastructure/cloud-regions/cloud-region/CloudOwner/RegionOne
   {
       "cloud-owner": "CloudOwner",
       "cloud-region-id": "RegionOne",
       "cloud-type": "SharedNode",
       "owner-defined-type": "OwnerType",
       "cloud-region-version": "v1",
       "identity-url": "http://10.12.25.2:5000/v2.0",
       "cloud-zone": "CloudZone",
       "resource-version": "1559336510793",
       "relationship-list": {
           ... ...

4. Add route on sdnc cluster VM node, which is the cluster VM node where pod sdnc-sdnc-0 is running on. This will allow ONAP SDNC to configure BRG later on. 
 
::

   ip route add 10.3.0.0/24 via 10.0.101.10 dev ens3

5. Initialize SDNC ip pool by running command from Rancher node 

:: 

   kubectl -n onap exec -it dev-sdnc-sdnc-0 -- /opt/sdnc/bin/addIpAddresses.sh VGW 10.5.0 22 250

6. Install Python and other Python libraries

::
 
   integration/test/vcpe/bin/setup.sh


7. Change the Openstack env parameters and one customer service related parameter in vcpecommon.py

:: 

    cloud = { 
        '--os-auth-url': 'http://10.12.25.2:5000',
        '--os-username': 'xxxxxxxxxx',
        '--os-user-domain-id': 'default',
        '--os-project-domain-id': 'default',
        '--os-tenant-id': 'xxxxxxxxxxxxxxxx' if oom_mode else '1e097c6713e74fd7ac8e4295e605ee1e',
        '--os-region-name': 'RegionOne',
        '--os-password': 'xxxxxxxxxxx',
        '--os-project-domain-name': 'xxxxxxxxx' if oom_mode else 'Integration-SB-07',
        '--os-identity-api-version': '3' 
    }   

    common_preload_config = { 
        'oam_onap_net': 'xxxxxxxx' if oom_mode else 'oam_onap_lAky',
        'oam_onap_subnet': 'xxxxxxxxxx' if oom_mode else 'oam_onap_lAky',
        'public_net': 'xxxxxxxxx',
        'public_net_id': 'xxxxxxxxxxxxx'
    }   

::

    # CHANGEME: vgw_VfModuleModelInvariantUuid is in rescust service csar, open service template with filename like service-VcpesvcRescust1118-template.yml and look for vfModuleModelInvariantUUID under groups vgw module metadata. 
    self.vgw_VfModuleModelInvariantUuid = 'xxxxxxxxxxxxxxx'

8. Initialize vcpe

::
   
   vcpe.py init

9. Run a command from Rancher node to insert vcpe customer service workflow entry in SO catalogdb. You should be able to see a sql command printed out from the above step output at the end, and use that sql command to replace the sample sql command below (inside the double quote) and run it from Rancher node:

::

   kubectl exec dev-mariadb-galera-mariadb-galera-0 -- mysql -uroot -psecretpassword catalogdb -e "INSERT INTO service_recipe (ACTION, VERSION_STR, DESCRIPTION, ORCHESTRATION_URI, SERVICE_PARAM_XSD, RECIPE_TIMEOUT, SERVICE_TIMEOUT_INTERIM, CREATION_TIMESTAMP, SERVICE_MODEL_UUID) VALUES ('createInstance','1','vCPEResCust 2019-06-03 _04ba','/mso/async/services/CreateVcpeResCustService',NULL,181,NULL, NOW(),'6c4a469d-ca2c-4b02-8cf1-bd02e9c5a7ce')"

10. Run Robot to create and distribute for vCPE customer service. This step assumes step 1 has successfully distributed all vcpe models except customer service model

::

   ete-k8s.sh onap distributevCPEResCust

11. Manually copy vCPE customer service csar (starting with service-Vcperescust) under Robot container /tmp/csar directory to Rancher vcpe/csar directory, now you should have these files:

::

    root@sb00-nfs:~/integration/test/vcpe/csar# ls -l
    total 528
    -rw-r--r-- 1 root root 126545 Jun 26 11:28 service-Demovcpeinfra-csar.csar
    -rw-r--r-- 1 root root  82053 Jun 26 11:28 service-Demovcpevbng-csar.csar
    -rw-r--r-- 1 root root  74179 Jun 26 11:28 service-Demovcpevbrgemu-csar.csar
    -rw-r--r-- 1 root root  79626 Jun 26 11:28 service-Demovcpevgmux-csar.csar
    -rw-r--r-- 1 root root  78156 Jun 26 11:28 service-Demovcpevgw-csar.csar
    -rw-r--r-- 1 root root  83892 Jun 26 11:28 service-Vcperescust20190625D996-csar.csar

12. Instantiate vCPE infra services

::

    vcpe.py infra

13. Install curl command on sdnc container inside sdnc-sdnc-0 pod

::

    sudo apk add curl

14. From Rancher node run vcpe healthcheck command to check connectivity from sdnc to brg and gmux, and vpp configuration of brg and gmux. Write down BRG MAC address printed out at the last line

::

    healthcheck-k8s.py onap

15. Instantiate vCPE customer service. Input the BRG MAC when prompt

::

    vcpe.py customer

16. Update libevel.so in vGMUX VM and restart the VM. This allows vGMUX to send events to VES collector in close loop test. See tutorial wiki for details

17. Run heatbridge. The heatbridge command usage: demo-k8s.sh <namespace> heatbridge <stack_name> <service_instance_id> <service> <oam-ip-address>, please refer to vCPE tutorial page on how to fill in those paraemters. See an example as following:

::

    ~/integration/test/vcpe# ~/oom/kubernetes/robot/demo-k8s.sh onap heatbridge vcpe_vfmodule_e2744f48729e4072b20b_201811262136 d8914ef3-3fdb-4401-adfe-823ee75dc604 vCPEvGMUX 10.0.101.21

18. Push vCPE closed loop Policy. Copy the two operational policy from vcpe/preload_templates to Robot container and then run the following two commands inside Robot container. You can find more details in JIRA INT-1089 - Create vCPE closed loop policy and push to policy engine

::

    curl -k --silent --user 'healthcheck:zb!XztG34' -X POST "https://policy-api:6969/policy/api/v1/policytypes/onap.policies.controlloop.Operational/versions/1.0.0/policies" -H "Accept: application/json" -H "Content-Type: application/json" -d @operational.vcpe.json.txt
    curl --silent -k --user 'healthcheck:zb!XztG34' -X POST "https://policy-pap:6969/policy/pap/v1/pdps/policies" -H "Accept: application/json" -H "Content-Type: application/json" -d @operational.vcpe.pap.json.txt

19. Start closed loop test by triggering packet drop VES event, and monitor if vGMUX is restarting. You may need to run the command twice if the first run fails

:: 

    vcpe.py loop


Test Status
~~~~~~~~~~~~~~~~~~~~~
The use case has been tested for Dublin release, the test report can be found on https://wiki.onap.org/display/DW/vCPE+%28Heat%29+-+Dublin+Test+Status

Known Issues and Workaround
~~~~~~~~~~~~~~~~~~~~~~~~~~~~
1) NATs are installed on BRG and vBNG. In order to allow SDNC to send BRG configuration message through vBNG, SDNC host VM IP address is preloaded on BRG and vBNG during VM instantiation, and provisioned into the NATs. If SDNC changes its host VM, SDNC host VM IP changes and we need to manually update the IP in /opt/config/sdnc_ip.txt. Then run:

::

  root>vppctl tap delete tap-0
  root>vppctl tap delete tap-1
  root>/opt/nat_service.sh
  root>vppctl restart

2) During vCPE customer service instantiation, though vGW should come up successfully BRG vxlan tunnel configuration is likely to fail in SDNC cluster environment due to SDNC unreachable to BRG. See more detail in JIRA INT-1127. One workaround is to run vCPE use case with SDNC cluster disabled.
