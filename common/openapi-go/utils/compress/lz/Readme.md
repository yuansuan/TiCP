
# FLOATING-POINT DATA COMPRESSION

1. **lz算法**

参考 [http://leobago.com/projects/lz/](http://leobago.com/projects/lz/)

It receives a block of floating-point numbers (one per row) and decomposes each floating-point number into an array of bytes, forming a matrix of bytes. The matrix has 4 or 8 columns (for single and double precision respectively) and as many rows as floating-point numbers in the block. Then, it transposes the matrix resulting in a new matrix where the first rows show low entropy (exponent) and the last rows have high entropy (last mantissa bits). Then, lz compresses the first rows, keep the high entropy rows uncompressed and if desired by the user, discard the last rows.

将float数组，按一个float一行，把每个byte看成一个格子，组成矩阵。如果数据是连续的，大概率每列上的值也是相近的。然后对每列用别的算法（zip、gzip之类的）进行压缩。

参考里的代码可能是边界条件没处理好，在小数据时会有问题。

性能对比

 ---   compress ratio(compared with origin file size):

  use rob.odb
| float | 254132624 | 61222900 | 76305923 | (lz)0.24090925059664908 | (gzip)0.30026024128252027 | (lz in c)0.23979692981094786|
| double | 254132624 | 61493407 | 76305923 | (lz)0.24197368300104594 | (gzip)0.30026024128252027 |

 use d3plot01
|  float | 14311168 | 20096271 | 23602880  |  (lz)0.17580321635765284 |  (gzip)0.20647921294969185 |
|  double | 114311168 | 20816458 23602880 | (lz)0.18210344941974524 | (gzip)0.20647921294969185 |


 ---   compress speed 
 use d3plot01
|  lz float | 114311168 | 2.00834519      |   (lz float) 54.281318541659665 |
|  lz double | 114311168 | 2.136632822 |     (lz double)51.02216154198908 |
|  gzip | 114311168  | 1.827797696    |        (gzip)     59.64315702912452 |

 use rob.odb
|  lz float | 254132624 | 4.253635564  |   (lz float) 56.97708532458719 |
|  lz double | 254132624 | 4.604696583 |      (lz double)52.633165313104534 |
|  gzip | 254132624 | 3.341635055    |        (gzip)     72.52729651225381 |