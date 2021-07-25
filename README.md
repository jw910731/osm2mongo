# osm2mongo
This tool dumps osm pbf data to mongodb in GeoJson format. Other attributes (or tags) is scattered in the document.
The only osm type it can currently handle is Node, corresponding to GeoJson Point type. Other objects are simple ignored.
This is a rough tool that is written in a sleepless night and is not reliable. Use at your own risk.

BTW, the go routine amount of parallel decoder are hard coded. 