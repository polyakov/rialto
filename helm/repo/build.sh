#!/bin/bash

cd charts

rm ldap.ApacheDS-0.1.0.tgz
tar -czvf ldap.ApacheDS-0.1.0.tgz ldap.ApacheDS/

rm openid-connect-ldap-mitre-0.1.0.tgz
tar -czvf openid-connect-ldap-mitre-0.1.0.tgz openid-connect-ldap-mitre/

cd ..
