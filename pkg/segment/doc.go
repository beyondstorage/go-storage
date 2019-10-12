/*
Package segment provided segment metadata support for storage.

This package handles segment operation details, and Storager implementer don't need to write their own segment logic.

Segment splits a File into several virtual parts, every part will have it's offset and size. Segment's life cycle
contains five stage.

                                                        +----------------+
                                                        |                |
                                                 +------>    Complete    |
                                                 |      |                |
                                                 |      +----------------+
                                                 |
       +------------+         +-------------+    |
       |            |         |             |    |
       |    Init    +-------->+    Write    +----+
       |            |         |             |    |
       +------------+         +-------------+    |
                                                 |
                                                 |      +-------------+
                                                 |      |             |
                                                 +------>    Abort    |
                                                        |             |
                                                        +-------------+
                              +------------+
                              |            |
                              |    Read    |
                              |            |
                              +------------+

At Init stage

New segment will be created with a ID, this ID could have different meanings for different servers.

At Write stage

New Part will be added into segment, and segment will check this Part has no intersected areas with existed Part.

At Complete stage

Segment will check every part is connected, make sure there is no hole for it.

At Abort stage

Segment will be destroyed, and no more Part could be added into this segment.

At Read stage

This stage is special and will not be handled by Segment, and don't need to call Init before use this API call.

Segment could be used in Storager segment API set
*/
package segment
