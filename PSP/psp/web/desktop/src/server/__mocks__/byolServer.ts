/* Copyright (C) 2016-present, Yuansuan.cn */

const data = [
  {
    app_id: 'app_id',
    app_name: 'sample',
    version: '2.0',
    active: true,
    merchandise_id: 'merchandise_id',
  },
  {
    app_id: 'app_id',
    app_name: 'sample',
    version: '2.0',
    active: true,
    merchandise_id: 'merchandise_id',
  },
]

const Item = {
  app_id: 'app_id',
  app_name: 'sample',
  version: '2.0',
  active: true,
  merchandise_id: 'merchandise_id',
}

const license = {
  id: 'license_id',
  license: JSON.stringify({
    ip: '12.32.23.43',
    license_port: 34,
    extra_port: 234,
    provider_port: 234,
  }),
}

export const byolServer = {
  __data__: {
    data,
    Item,
    license,
  },
  get: async () => ({ data }),
  getlicense: async () => ({
    data: { id: license.id, license: license.license },
  }),
}
