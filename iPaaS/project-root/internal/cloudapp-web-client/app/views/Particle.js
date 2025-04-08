import P from 'react-tsparticles'
import React, { PureComponent } from 'react'

const types = {
  lines: {
    particles: {
      opacity: {
        value: 0.2
      },
      number: {
        value: 40
      },
      move: {
        speed: 1
      },
      line_linked: {
        enable: true,
        opacity: 0.3
      },
      size: {
        value: 3
      }
    }
  },
  nightSky: {
    particles: {
      number: {
        value: 200,
        density: {
          enable: true,
          value_area: 1000
        }
      },
      line_linked: {
        enable: true,
        opacity: 0.1
      },
      move: {
        direction: 'right',
        speed: 0.1
      },
      size: {
        value: 1
      },
      opacity: {
        anim: {
          enable: true,
          speed: 1,
          opacity_min: 0.1
        }
      }
    },

    retina_detect: true
  },
  snow: {
    particles: {
      opacity: {
        anim: {
          enable: true,
          speed: 2,
          opacity_min: 0.05
        },
        opacity_max: 0.1
      },
      number: {
        value: 40,
        density: {
          enable: false
        }
      },
      size: {
        value: 8,
        random: true
      },
      move: {
        direction: 'bottom',
        out_mode: 'out'
      },
      line_linked: {
        enable: false
      }
    }
  },
  bubble: {
    particles: {
      number: {
        value: 160,
        density: {
          enable: false
        }
      },
      size: {
        value: 3,
        random: true,
        anim: {
          speed: 4,
          size_min: 0.3
        }
      },
      line_linked: {
        enable: false
      },
      move: {
        random: true,
        speed: 1,
        direction: 'top',
        out_mode: 'out'
      }
    },
    interactivity: {
      events: {
        onclick: {
          enable: true,
          mode: 'repulse'
        }
      },
      modes: {
        bubble: {
          distance: 250,
          duration: 2,
          size: 0,
          opacity: 0
        },
        repulse: {
          distance: 200,
          duration: 4
        }
      }
    }
  }
}
export default class Particles extends PureComponent {
  render() {
    const particlesInit = main => {
      console.log(main)
    }

    const particlesLoaded = container => {
      console.log(container)
    }
    const { type = 'snow' } = this.props
    return (
      <P
        id="tsparticles"
        init={particlesInit}
        loaded={particlesLoaded}
        options={{
          background: {
            color: {
              value: '#0d47a1'
            }
          },
          fpsLimit: 60,
          detectRetina: true
        }}
      />
    )
  }
}
