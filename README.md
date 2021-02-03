# grafana-stock-tracker

This repository will hold the application code for attempting to create a free stock tracker using the yahoo finance api. ( using https://github.com/achannarasappa/ticker as an influence ) 


The web server will act as a prometheus exporter. It will read a file and check for a sticker update
And if there's a new sticker it will start collecting stock price for that sticker and emit a prometheus metrics. 
Currently the plan is to leverage prometheus push gateway to publish stock metrics
