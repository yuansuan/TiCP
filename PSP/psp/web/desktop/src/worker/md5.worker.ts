import SparkMD5 from 'spark-md5'

const _eventName = 'md5'

const ctx: Worker = self as any

function md5(chunkTotal, CHUNK_SIZE, file, uploadId) {
  const spark = new SparkMD5.ArrayBuffer()

  let blobSlice = File.prototype.slice,
    currentChunk = 0
  const fileReader = new FileReader()

  fileReader.onload = e => {
    spark.append(e.target.result) // Append array buffer
    currentChunk++

    if (currentChunk < chunkTotal) {
      loadNext()
    } else {
      ctx.postMessage({
        eventName: _eventName,
        eventData: { md5: spark.end(), uploadId },
      })
    }
  }

  fileReader.onerror = function() {
    ctx.postMessage({
      eventName: _eventName,
      eventData: { md5: 'md5计算失败', uploadId },
    })
  }

  function loadNext() {
    let start = currentChunk * CHUNK_SIZE,
      end = start + CHUNK_SIZE >= file.size ? file.size : start + CHUNK_SIZE

    fileReader.readAsArrayBuffer(blobSlice.call(file, start, end))
  }

  loadNext()
}

ctx.addEventListener('message', event => {
  // 计算md5
  const { eventName, eventData } = event.data
  if (eventName === _eventName) {
    const { chunkTotal, CHUNK_SIZE, file, uploadId } = eventData
    md5(chunkTotal, CHUNK_SIZE, file, uploadId)
  }
})

export default null as any
