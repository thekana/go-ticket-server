#!bin/sh
mkdir /etc/nginx/sites-enabled
ln -s /etc/nginx/sites-available/ticketeer.ml.conf /etc/nginx/sites-enabled/ticketeer.ml.conf
ln -s /etc/nginx/sites-available/goticket.tk.conf /etc/nginx/sites-enabled/goticket.tk.conf