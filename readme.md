convert input to 64k mono aac output
1. ffmpeg -i input.file -ac 1 -c:a aac -b:a 64k output.aac

should recompile ffmpeg with libfdk_aac to support streaming

