## pixel_boot_img
automaton download and extract `init_boot.img` from pixel factory image,
which used by magisk.

## usage
to trigger a download and processing for a factory image, `open(reopen)` a issue which has a factory image link download link as last line of issue body.

e.g.
> *(can be found at https://developers.google.com/android/images)  
> https://dl.google.com/dl/android/aosp/rango-bd3a.250808.001-factory-8f7b23b7.zip

if the github action is finished successfully. result will be committed at the `storage-[device]-[build]` dirs.

e.g.
```
storage/rango/BD3A.250808.001/
```