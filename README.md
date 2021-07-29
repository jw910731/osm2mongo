# osm2mongo
This tool dumps osm pbf data to mongodb in GeoJson format. Other attributes (or tags) is scattered in the document.
The only osm type it can currently handle is Node, corresponding to GeoJson Point type. Other objects are simple ignored.
This is a rough tool that is written in a sleepless night and is not reliable. Use at your own risk.

BTW, the goroutine amount of parallel decoder are hard coded, change it to fit your needs. 

# Schema
The schema of imported data is:
```json
{
        "_id": "OSM ID (int)",
        "coordinates": ["Lontitude (number)", "Latitude (number)"],
        "type": "Point", // this is static value
        "time": "OSM info timestamp (timestamp)",
        //... Other OSM tags
}
```
If the tag look like "xxx:yyy", "xxx:zzz". It will be gather to a document like `xxx: {yyy: value, zzz: value2}`

# License
```
        DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE 
                    Version 2, December 2004 

 Copyright (C) 2004 Sam Hocevar <sam@hocevar.net> 

 Everyone is permitted to copy and distribute verbatim or modified 
 copies of this license document, and changing it is allowed as long 
 as the name is changed. 

            DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE 
   TERMS AND CONDITIONS FOR COPYING, DISTRIBUTION AND MODIFICATION 

  0. You just DO WHAT THE FUCK YOU WANT TO.
  ```