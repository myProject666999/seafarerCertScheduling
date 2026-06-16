import http from './index'

export const getSeafarers = (params) => http.get('/seafarers/', { params })
export const getSeafarer = (id) => http.get(`/seafarers/${id}`)
export const createSeafarer = (data) => http.post('/seafarers/', data)
export const updateSeafarer = (id, data) => http.put(`/seafarers/${id}`, data)
export const deleteSeafarer = (id) => http.delete(`/seafarers/${id}`)

export const getCertTypes = () => http.get('/certificate-types/')
export const createCertType = (data) => http.post('/certificate-types/', data)
export const deleteCertType = (id) => http.delete(`/certificate-types/${id}`)

export const getSeafarerCerts = (params) => http.get('/seafarer-certificates/', { params })
export const createSeafarerCert = (data) => http.post('/seafarer-certificates/', data)
export const updateSeafarerCert = (id, data) => http.put(`/seafarer-certificates/${id}`, data)
export const deleteSeafarerCert = (id) => http.delete(`/seafarer-certificates/${id}`)

export const getShips = (params) => http.get('/ships/', { params })
export const getShip = (id) => http.get(`/ships/${id}`)
export const createShip = (data) => http.post('/ships/', data)
export const updateShip = (id, data) => http.put(`/ships/${id}`, data)
export const deleteShip = (id) => http.delete(`/ships/${id}`)

export const getShipPositions = (params) => http.get('/ship-positions/', { params })
export const createShipPosition = (data) => http.post('/ship-positions/', data)
export const deleteShipPosition = (id) => http.delete(`/ship-positions/${id}`)

export const getAssignments = (params) => http.get('/assignments/', { params })
export const createAssignment = (data) => http.post('/assignments/', data)
export const disembarkAssignment = (id, data) => http.post(`/assignments/${id}/disembark`, data)

export const getTransfers = (params) => http.get('/transfers/', { params })
export const createTransfer = (data) => http.post('/transfers/', data)
export const approveTransfer = (id, data) => http.post(`/transfers/${id}/approve`, data)
export const rejectTransfer = (id, data) => http.post(`/transfers/${id}/reject`, data)
export const cancelTransfer = (id) => http.post(`/transfers/${id}/cancel`)

export const getAlerts = (params) => http.get('/alerts/', { params })
export const getAlertStats = () => http.get('/alerts/stats')
export const handleAlert = (id, data) => http.post(`/alerts/${id}/handle`, data)
export const runAlertScan = () => http.post('/alerts/scan')

export const getContracts = (params) => http.get('/contracts/', { params })
export const createContract = (data) => http.post('/contracts/', data)

export const getEmbarkRecords = (params) => http.get('/embark-records/', { params })
export const createEmbarkRecord = (data) => http.post('/embark-records/', data)

export const getLeaveRecords = (params) => http.get('/leave-records/', { params })
export const createLeaveRecord = (data) => http.post('/leave-records/', data)
export const endLeave = (id) => http.post(`/leave-records/${id}/end`)

export const getHealthRecords = (params) => http.get('/health-records/', { params })
export const createHealthRecord = (data) => http.post('/health-records/', data)
