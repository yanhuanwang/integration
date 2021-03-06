#! /usr/bin/python

import time
import logging
import json
import mysql.connector
import ipaddress
import re
import sys
import base64
from vcpecommon import *
import preload
import commands
import vcpe_custom_service


logging.basicConfig(level=logging.INFO, format='%(message)s')

cpecommon = VcpeCommon()
custom = vcpe_custom_service.CustomService(cpecommon)

nodes=['mux']
hosts = cpecommon.get_vm_ip(nodes)

custom.del_vgmux_ves_mode(hosts['mux'])
time.sleep(2)
custom.del_vgmux_ves_collector(hosts['mux'])
#exit()

time.sleep(2)
logging.info('Setting vGMUX DCAE collector IP address')
custom.set_vgmux_ves_collector(hosts['mux'])
time.sleep(2)
vgmux_vnf_name = cpecommon.load_object('vgmux_vnf_name')
logging.info('vGMUX VNF instance name is %s', vgmux_vnf_name)
logging.info('Letting vGMUX report packet loss to DCAE')
custom.set_vgmux_packet_loss_rate(hosts['mux'], 55, vgmux_vnf_name)
