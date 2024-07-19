# CDN With Aliyun OSS Storage (preview)
> This plugin can be used to store static files to Aliyun OSS.

## How to use

### Build
```bash
./answer build --with github.com/apache/incubator-answer-plugins/cdn-aliyun
```

### Configuration
- `Endpoint` -  Endpoint of AliCloud OSS storage, such as oss-cn-hangzhou.aliyuncs.com
- `Bucket Name` - Your bucket name
- `Object Key Prefix` - Prefix of the object key like 'static/' that ending with '/'
- `Access Key Id` - AccessKeyID of the AliCloud OSS storage
- `Access Key Secret` - AccessKeySecret of the AliCloud OSS storage
- `Visit Url Prefix` - Prefix of access address for the CDN file, ending with '/' such as https://static.example.com/xxx/
- `Max File Size` - Max file size in MB, default is 10MB