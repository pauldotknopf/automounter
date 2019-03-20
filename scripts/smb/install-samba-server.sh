#!/usr/bin/env bash
set -e
sudo apt-get install -y samba cifs-utils
sudo smbpasswd -a pknopf
sudo mkdir -p /var/lib/smb-custom
sudo chmod 777 /var/lib/smb-custom
echo "[custom]" | sudo tee -a /etc/samba/smb.conf
echo "  path = /var/lib/smb-custom" | sudo tee -a /etc/samba/smb.conf
echo "  valid users = pknopf" | sudo tee -a /etc/samba/smb.conf
echo "  read only = no" | sudo tee -a /etc/samba/smb.conf
sudo systemctl restart smbd
