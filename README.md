# olx.ph scraper

Searches olx.ph for pattern and stores matching items in json
file.  Captures the following information about an item:
- short description
- price (PHP)
- link to item

Supports use of price ranges to narrow down results.  Use ```-a```
and ```-b``` to specify the minimum and maximum price (PHP), 
respectively.

Use ```-p``` to print the new items found from last run.

Use ```-j``` to specify a custom location for the json file.

Public docker container can be found on docker hub under
 ```pehowell/olx-scraper```. When running the docker container,
you should specify a mounted volume for /data (or your custom
 location for the json file) so the json can persist between runs.

Docker example:

```docker run --rm -v /my/home/dir:/data pehowell/olx-scraper -p -a 1000 -b 2000 gameboy```

This searches for "gameboy" between the price of 1000 and 2000 PHP
and prints any new items found between runs.  The ```olx.json``` file
will be stored in ```/my/home/dir/```.