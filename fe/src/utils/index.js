import { version } from '../../package.json'

const PackageVersion = version

// Units are K,M,G,T,P (powers of 1024) like `/usr/bin/ls`.
const calcFromBytes = (raw) => {
  if (typeof raw === 'string') {
    raw = parseInt(raw)
  }
  if (raw >= 1125899906842624) {
    return (raw / 1125899906842624).toFixed(1) + ' P'
  } else if (raw >= 1099511627776) {
    return (raw / 1099511627776).toFixed(1) + ' T'
  } else if (raw > 1073741824) {
    return (raw / 1073741824).toFixed(1) + ' G'
  } else if (raw > 1048576) {
    return (raw / 1048576).toFixed(1) + ' M'
  } else if (raw > 1024) {
    return (raw / 1024).toFixed(1) + ' K'
  } else {
    return raw.toFixed(1) + ' B'
  }
}

export {
  PackageVersion,
  calcFromBytes
}
