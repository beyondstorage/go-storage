---
author: Xuanwo <github@xuanwo.io>
status: finished
updated_at: 2020-01-13
---

# Proposal: Normalize metadata storage class

## Background

> Simple introduce about `storage class`
>
> Storage services usually offer different storage class for different use cases.
> - For hot data which frequently accessed, service does optimization for low latency and high throughput performance. 
> - For cold data which accessed one or two times over years, service does optimization for durability and low cost.
>
> Services will mark different storage class in object metadata.

As described in proposal [6-normalize-metadata](https://github.com/Xuanwo/storage/blob/master/docs/design/6-normalize-metadata.md), storage class is a meta that has not been standardized. We need to normalize it so that we can use them across different services.

We have following problems to address:

### Different name

Different services may have different name. 

As a private metadata, services may assign them in private headers:

- `qingstor` use `x-qs-storage-class` to carry the `storage-class`
- `s3` use `x-amz-storage-class` to carry the `storage-class`

Even storage class itself is not a common usage:

- `azblob` use `AccessTier` to represent different storage class

### Different value

It's obvious different provider won't use the same meaning value.

- `s3` uses `Standard`, `Standard-IA`, `Intelligent-Tiering`, `Glacier` and so on.
- `azblob` uses `Hot`, `Cool` and `Archive`.
- `gcs` uses `Standard`, `Nearline` and `Coldline`.
- ...

Almost every provider has its own set of storage class values.

### Different TTFB

> TTFB: Time to first byte

Different storage class means different under layer storage media or algorithm which provides different warrant about TTFB.

- `s3`: `Standard` TTFB maybe milliseconds while `Glacier` needs `select minutes or hours`
- `gcs`: All storage class provides the same latency
- `kodo`: `Archive` needs another API call to restore the files
- `cos`: `Archvie` needs to submit application to restore the files
- ...

We need to overcome those differences.

## Proposal

So I propose following changes:

### Normalize name

Metadata for `storage class`'s name will be normalized to `storage-class`.

### Add global storage class set

Add `Hot`, `Warm`, `Cold` storage classes as a global storage class set across storage services.

`Hot` is used for frequently accessed data.
 
 - All service supports this kind of storage class
 - All services' default storage class
 - If a service doesn't have an idea for storage class, it provides `Hot` storage class
 
`Warm` is used for infrequent access data, maybe accessed several times a month.
 
 - Higher latency and lower performance compared to `Hot`
 - `Read` operation requires extra fees (depend on services)
 
`Cold` is used for archiving data which maybe accessed one or two times a year.

- `Cold` only support write operation
- Depends on services' implementations, `Cold`, a.k.a, `Archive` storage class may need extra time (minutes to hours, except `gcs`) or extra API (`kodo`) even manual application (`cos`). So we will not add extra support to read `Cold` data, just return the error.

In conclusion:

- `Hot` will be mapped to service's fastest/default storage class
- `Warm` will be infrequent access storage class for service
- `Cold` will be the archive storage class for service

For example:

- s3 `Standard` == `Hot`
- s3 `Standard` == `Warm`
- s3 `Glacier` == `Cold`

### Service behavior spec

- Read storage class from pairs, service need to parse from global storage class
- Write storage class to metadata, service need to format into global storage class
- Meet unsupported storage class while parsing from global storage class, return `ErrStorageClassNotSupported` with service name and storage class name is services have more than one storage class
- Meet unsupported storage class while formatting into global storage class, return `ErrStorageClassNotSupported` with service name and storage class name

## Rationale

### Services storage classes

#### azblob

ref: https://docs.microsoft.com/en-us/azure/storage/blobs/storage-blob-storage-tiers

- Hot: Optimized for storing data that is accessed frequently.
- Cool: Optimized for storing data that is infrequently accessed and stored for at least 30 days.
- Archive: Optimized for storing data that is rarely accessed and stored for at least 180 days with flexible latency requirements (on the order of hours).

#### cos

ref: https://cloud.tencent.com/document/product/436/33417

- Standard: COS Standard is an object storage service with high reliability, availability, and performance.Its low latency and high throughput make it well suitable for the use cases involving lots of hotspot files or frequent data access.
- Standard_IA: COS Infrequent Access (COS Standard_IA) provides object storage services featuring high reliability and low storage cost and access latency. It offers lowered pricing for storage and keeps the first-byte access time within milliseconds, ensuring that data can be fast retrieved with no wait required. However, data retrieval incurs fees. It is suitable for business scenarios where the access frequency is low (e.g., once or twice per month).
- Archive: Archive Storage is a highly reliable object storage service that has ultra-low storage cost and long-term data retention. Featuring the lowest storage price, Archive Storage needs a longer time to read data and is suitable for archived data that needs to be stored for a long time.

#### dropbox

None

#### fs

None

#### gcs

ref: https://cloud.google.com/storage/docs/storage-classes

- Standard: Standard Storage is best for data that is frequently accessed ("hot" data) and/or stored for only brief periods of time.
- Nearline: Nearline Storage is a low-cost, highly durable storage service for storing infrequently accessed data. Nearline Storage is a better choice than Standard Storage in scenarios where slightly lower availability, a 30-day minimum storage duration, and costs for data access are acceptable trade-offs for lowered at-rest storage costs.
- Coldline: Coldline Storage is a very-low-cost, highly durable storage service for storing infrequently accessed data. Coldline Storage is a better choice than Standard Storage or Nearline Storage in scenarios where slightly lower availability, a 90-day minimum storage duration, and higher costs for data access are acceptable trade-offs for lowered at-rest storage costs.
- Archive: Archive Storage is the lowest-cost, highly durable storage service for data archiving, online backup, and disaster recovery. Unlike the "coldest" storage services offered by other Cloud providers, your data is available within milliseconds, not hours or days.

#### kodo

ref: https://developer.qiniu.com/kodo/api/3710/chtype

- STANDARD (0, 标准存储)
- STANDARD_IA (1, 低频存储)
- ARCHIVE (2, 归档存储)

#### oss

ref: https://www.alibabacloud.com/help/doc-detail/51374.htm

- Standard: OSS Standard storage provides highly reliable, highly available, and high-performance object storage services that support frequent data access. The high-throughput and low-latency service response capability of OSS can effectively support access to hot data. Standard storage is ideal for storing images for social networking and sharing, and storing data for audio and video applications, large websites, and big data analytics. 
- IA: OSS IA storage is suitable for storing long-lived, but less frequently accessed data (an average of once or twice per month). IA storage offers a storage unit price lower than that of Standard storage, and is suitable for long-term backup of various mobile apps, smart device data, and enterprise data. It also supports real-time data access. Objects of the IA storage class have a minimum storage period.
- Archive: OSS Archive storage has the lowest price among the three storage classes. It is suitable for long-term (at least half a year) storage of data that is infrequently accessed. The data may take about one minute to restore before it can be read. This storage option is suitable for data such as archival data, medical images, scientific materials, and video footage.

#### qingstor

ref: https://www.qingcloud.com/products/qingstor/

- Standard
- Standard_IA

#### s3

ref: https://docs.aws.amazon.com/AmazonS3/latest/dev/storage-class-intro.html

- STANDARD: Frequently accessed data
- STANDARD_IA: Long-lived, infrequently accessed data
- INTELLIGENT_TIERING: Long-lived data with changing or unknown access patterns
- ONEZONE_IA: Long-lived, infrequently accessed, non-critical data
- GLACIER: 	Long-term data archiving with retrieval times ranging from minutes to hours
- DEEP_ARCHIVE: Archiving rarely accessed data with a default retrieval time of 12 hours

#### uss

Not GA

### Why `Hot`, `Warm` and `Cold`?

Service provider always wraps their product into fancy names, we need to figure out the core value. Instead of answering `Why`, I prefer `Why not`.

#### Why not `Standard`?

Yes, `Standard` are both widely used, 6/7 services picked `Standard` for their default storage class. 

However, I don't know what `Standard` for. As a newbie without any idea of storage class, how can I realize `Standard` is for frequently accessed data in 10 seconds? What's `Standard`? Are there any non-`Standard`?

From my point of view, others picked this name just because `S3` picked it. It's a `Standard` for `Amazon Web Services`.

#### Why not `Nearline`?

`Nearline` has a long history which dates back to [IBM 3850 Mass Storage System (MSS) tape library in 1974](https://en.wikipedia.org/wiki/Nearline_storage). For developers in the 2020s, fewer and fewer of them know about this idea except for those who focus on storage.
  
As an application-oriented unified storage layer, it makes more sense to let users choose their data type instead of storage type.
 
In conclusion:
 
`Hot`, `Warm` and `Cold` are clear, no confusion, understandable.

Needless to say, they are much shorter :-)
 
## Compatibility

All API call with storage class will be affected.

## Implementation

Most of the work would be done by the author of this proposal.