# olx.ph scraper

Searches olx.ph for pattern and stores matching items in json
file.  Captures the following information about an item:
- short description
- price
- link to item

Supports use of price ranges to narrow down results.  Use ```-a```
and ```-b``` to specify the minimum and maximum price, 
respectively.

Use ```-p``` to print the new items found from last run.