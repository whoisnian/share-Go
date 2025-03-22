const FETCH_OPTIONS = {
  credentials: 'same-origin',
  headers: {
    accept: 'application/json',
    'content-type': 'application/json'
  }
}

const fetchGetHead = async (url) => {
  const result = await window.fetch(url, { method: 'GET', ...FETCH_OPTIONS })
  if (!result.ok) throw new Error(`[fetchGetHead] failed with ${result.status}: ${url}`)
}

const fetchPostHead = async (url, data) => {
  const result = await window.fetch(url, { method: 'POST', body: JSON.stringify(data), ...FETCH_OPTIONS })
  if (!result.ok) throw new Error(`[fetchPostHead] failed with ${result.status}: ${url}`)
}

const fetchDeleteHead = async (url, data) => {
  const result = await window.fetch(url, { method: 'DELETE', body: JSON.stringify(data), ...FETCH_OPTIONS })
  if (!result.ok) throw new Error(`[fetchDeleteHead] failed with ${result.status}: ${url}`)
}

const fetchPutHead = async (url, data) => {
  const result = await window.fetch(url, { method: 'PUT', body: JSON.stringify(data), ...FETCH_OPTIONS })
  if (!result.ok) throw new Error(`[fetchPutHead] failed with ${result.status}: ${url}`)
}

const fetchGetJSON = async (url) => {
  const result = await window.fetch(url, { method: 'GET', ...FETCH_OPTIONS })
  if (!result.ok) throw new Error(`[fetchPostJSON] failed with ${result.status}: ${url}`)
  return result.json()
}

const fetchPostJSON = async (url, data) => {
  const result = await window.fetch(url, { method: 'POST', body: JSON.stringify(data), ...FETCH_OPTIONS })
  if (!result.ok) throw new Error(`[fetchPostJSON] failed with ${result.status}: ${url}`)
  return result.json()
}

const fetchDeleteJSON = async (url, data) => {
  const result = await window.fetch(url, { method: 'DELETE', body: JSON.stringify(data), ...FETCH_OPTIONS })
  if (!result.ok) throw new Error(`[fetchDeleteJSON] failed with ${result.status}: ${url}`)
  return result.json()
}

const fetchPutJSON = async (url, data) => {
  const result = await window.fetch(url, { method: 'PUT', body: JSON.stringify(data), ...FETCH_OPTIONS })
  if (!result.ok) throw new Error(`[fetchPutJSON] failed with ${result.status}: ${url}`)
  return result.json()
}

const fetchGetJSONWithStatus = async (url) => {
  const result = await window.fetch(url, { method: 'GET', ...FETCH_OPTIONS })
  return {
    ok: result.ok,
    status: result.status,
    content: result.headers.get('content-type')?.includes('application/json') ? await result.json() : await result.text()
  }
}

const fetchPostJSONWithStatus = async (url, data) => {
  const result = await window.fetch(url, { method: 'POST', body: JSON.stringify(data), ...FETCH_OPTIONS })
  return {
    ok: result.ok,
    status: result.status,
    content: result.headers.get('content-type')?.includes('application/json') ? await result.json() : await result.text()
  }
}

const fetchDeleteJSONWithStatus = async (url, data) => {
  const result = await window.fetch(url, { method: 'DELETE', body: JSON.stringify(data), ...FETCH_OPTIONS })
  return {
    ok: result.ok,
    status: result.status,
    content: result.headers.get('content-type')?.includes('application/json') ? await result.json() : await result.text()
  }
}

const fetchPutJSONWithStatus = async (url, data) => {
  const result = await window.fetch(url, { method: 'PUT', body: JSON.stringify(data), ...FETCH_OPTIONS })
  return {
    ok: result.ok,
    status: result.status,
    content: result.headers.get('content-type')?.includes('application/json') ? await result.json() : await result.text()
  }
}

export {
  fetchGetHead,
  fetchPostHead,
  fetchDeleteHead,
  fetchPutHead,
  fetchGetJSON,
  fetchPostJSON,
  fetchDeleteJSON,
  fetchPutJSON,
  fetchGetJSONWithStatus,
  fetchPostJSONWithStatus,
  fetchDeleteJSONWithStatus,
  fetchPutJSONWithStatus
}
