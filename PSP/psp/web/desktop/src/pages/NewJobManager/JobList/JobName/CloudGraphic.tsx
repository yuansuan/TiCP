import React, { useState, useRef, useEffect } from 'react'
import { Button } from 'antd'
import Slider from 'react-slick'
import 'slick-carousel/slick/slick.css'
import 'slick-carousel/slick/slick-theme.css'
import { Http as JobLogHttp } from '@/utils/JobLogHttp'
import { Http as SnapshotHttp } from '@/utils/SnapshotHttp'
import debounce from 'lodash/debounce'
import {
  LeftCircleOutlined,
  RightCircleOutlined,
  PauseCircleOutlined,
  PlayCircleOutlined
} from '@ant-design/icons'
import styled from 'styled-components'
import { previewImage } from '@/components'

const StyledLayout = styled.div`
  .action {
    display: flex;
    align-content: center;
    justify-content: center;
    align-items: center;
    margin: 5px !important;
    .btns {
      .btn {
        display: none;
        margin: 0 5px;
      }
    }
  }
  .nav {
    display: flex;
    align-content: center;
    justify-content: center;
    align-items: center;

    .slidesNav {
      margin: 5px;
      width: calc(100% - 100px);

      .ant-carousel .slick-list .slick-slide.slick-current {
        background: #2196f3 !important;
      }

      .slick-list .slick-slide.slick-active.slick-current {
        background: #2196f3 !important;
      }
    }
  }
`

function SlideImg({ fileName, src, style, imgHeight, isStandardCompute }) {
  const imgRef = useRef(null)
  const [loading, setLoading] = useState(false)
  const [err, setErr] = useState(false)

  useEffect(() => {
    handleImageLoaded()
  }, [src])

  const handleImageLoaded = async () => {
    try {
      setLoading(true)
      const res = isStandardCompute ? await SnapshotHttp.get(src) : await JobLogHttp.get(src)
      imgRef.current.src = res.data?.snapshot
    } catch (e) {
      setErr(true)
    } finally {
      setLoading(false)
    }
  }

  return err ? (
    <div style={style}>图片加载出错</div>
  ) : (
    <div style={style}>
      <img
        ref={imgRef}
        height={imgHeight}
        onClick={() =>
          fileName && previewImage({ fileName, src: imgRef.current.src })
        }
      />
    </div>
  )
}

const getSlideFileName = slideSrc => {
  const urlSearchParams = new URLSearchParams(slideSrc)
  const fileName = urlSearchParams.get('path')
  return fileName.substr(fileName.lastIndexOf('/') + 1)
}

export const CloudGraphic = ({ data, isStandardCompute=false}) => {
  const slidesRef = useRef(null)
  const navRef = useRef(null)
  const playIntervalRef = useRef(null)
  const [slides, setSlides] = useState(null)
  const [nav, setNav] = useState(null)
  const [disabled, setDisabled] = useState(false)
  const [isPlaying, setIsPlaying] = useState(false)
  const [current, setCurrent] = useState(0)

  const goToLastSlide = () => {
    setCurrent(data.length - 1)
    slidesRef.current.slickGoTo(data.length - 1)
  }

  useEffect(() => {
    setSlides(slidesRef.current)
    setNav(navRef.current)

    goToLastSlide()

    return () => {
      playIntervalRef.current && clearInterval(playIntervalRef.current)
    }
  }, [data.length])

  const next = () => {
    if (current < data.length) navRef.current.slickNext()
  }

  const prev = () => {
    if (current > 0) navRef.current.slickPrev()
  }

  const play = () => {
    if (playIntervalRef.current) return

    setIsPlaying(true)
    playIntervalRef.current = setInterval(() => {
      slidesRef.current.next()
    }, 5000) as any
  }

  const stop = () => {
    if (playIntervalRef.current) clearInterval(playIntervalRef.current)
    setIsPlaying(false)
  }

  const contentStyle = {
    padding: '5px',
    height: '410px',
    backgroundColor: '#eee',
    display: 'flex',
    justifyContent: 'center'
  }

  const slidesSetting = {
    infinite: false,
    dots: false,
    lazyLoad: true,
    slidesToShow: 1,
    slidesToScroll: 1,
    autoplay: false,
    speed: 1000,
    arrows: true,
    autoplaySpeed: 5000
  }

  const navsContentStyle = {
    padding: '2px',
    height: '64px',
    backgroundColor: 'rgb(238, 238, 238, 0.5)',
    display: 'flex',
    justifyContent: 'center'
  }

  const navsSetting = {
    infinite: false,
    dots: false,
    slidesToShow: 6,
    slidesToScroll: 5,
    lazyLoad: true,
    swipeToSlide: true,
    focusOnSelect: true,
    speed: 600
  }

  return (
    <StyledLayout>
      <Slider
        {...(slidesSetting as any)}
        asNavFor={nav}
        ref={slidesRef}
        afterChange={current => {
          setCurrent(current)
          setDisabled(false)
        }}>
        {data.map(slide => {
          const fileName = getSlideFileName(slide)
          return (
            <div key={fileName}>
              <h3
                className='fileName'
                style={{ textAlign: 'center', height: 24 }}>
                {fileName}
              </h3>
              <SlideImg
                isStandardCompute={isStandardCompute}
                fileName={fileName}
                src={slide}
                style={contentStyle}
                imgHeight={'400px'}
              />
            </div>
          )
        })}
      </Slider>
      <div className='action'>
        <div className='btns'>
          <Button
            icon={isPlaying ? <PauseCircleOutlined /> : <PlayCircleOutlined />}
            disabled={disabled}
            shape='circle'
            className='btn'
            onClick={() => {
              isPlaying ? stop() : play()
            }}
            type={'primary'}
            title={isPlaying ? '暂停' : '播放'}
          />
          {`${current + 1}/${data.length}`}
        </div>
      </div>
      <div className='nav'>
        <Button
          icon={<LeftCircleOutlined />}
          disabled={disabled || current <= 0}
          style={{ zIndex: 1 }}
          shape='circle'
          onClick={debounce(() => {
            setDisabled(true)
            prev()
          }, 650)}
          title={'上一面板'}
        />
        <div className='slidesNav'>
          <Slider
            {...(navsSetting as any)}
            asNavFor={slides}
            ref={navRef}
            afterChange={() => {
              setDisabled(false)
            }}>
            {data.map(slide => {
              const fileName = getSlideFileName(slide)
              return (
                <div key={fileName} title={fileName}>
                  <SlideImg
                    isStandardCompute={isStandardCompute}
                    fileName={null}
                    src={slide}
                    style={navsContentStyle}
                    imgHeight={'60px'}
                  />
                </div>
              )
            })}
          </Slider>
        </div>
        <Button
          icon={<RightCircleOutlined />}
          shape='circle'
          style={{ zIndex: 1 }}
          disabled={disabled || current >= data.length - 1}
          onClick={debounce(() => {
            setDisabled(true)
            next()
          }, 650)}
          title={'下一面板'}
        />
      </div>
    </StyledLayout>
  )
}
