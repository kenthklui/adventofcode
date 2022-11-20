High level stategy for task 1 where we stitch together puzzle pieces:
* Read all tiles
* For each tile, compute a list of the 4 borders
* For each border, treat "." and "#" as 0s and 1s respectively, and compute a 10 bit unsigned integer used as a signature for each border (in retrospect, the string itself would have done fine; creating a binary signature was unnecessary)
* Compute all variants of each signature list, under all orientations: 4 rotations, multiplied by 2 flips, for 8 total variants per tile
* Create an exhaustive lookup map that allows searching any variant by the border signature integer, with the corresponding side (top/left/bottom/right)
* Using the above, do a backtracking-search to fill in each tile in the puzzle. Filter candidates by its neighbors; this will either be "border that matches the right side of the tile on the left", or "border that matches the bottom side of the tile above", or both.
* When all tiles are filled in, we are done.

Task 2:
* Stitch together the canvas using the previously computed solution. This required implementing additional tile rotation/flipping capabilities
* Instead of rotating the image, just rotate and flip the monster, and sum up the number of times it is found in each orientation. This assumes that monsters in the image don't overlap.
* Instead of highlighting which "#" belong to a monster versus not belonging to a monster, just count the number of total "#"s in the image and subtract `((number of monsters) * ("#"s per monster))`. Again, this assumes monsters in the image don't overlap.
