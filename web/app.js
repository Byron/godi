const form = document.querySelector('#controls');
const output = document.querySelector('#output');
const command = document.querySelector('#command');
const clientID = `${Date.now().toString(16)}${Math.random().toString(16).slice(2)}`;
const headers = {'Content-Type':'application/json','Client-ID':clientID};
let state;

const lines = name => form.elements[name].value.split('\n').map(v => v.trim()).filter(Boolean);
const readForm = () => ({
  ...state,
  mode: form.elements.mode.value,
  verbosity: form.elements.verbosity.value,
  format: form.elements.format.value,
  spid: Number(form.elements.spid.value),
  spod: Number(form.elements.spod.value),
  fep: form.elements.fep.value.split(',').map(v => v.trim()).filter(Boolean),
  sources: lines('sources'),
  destinations: lines('destinations'),
  verify: form.elements.verify.checked ? 'verify' : ''
});
const renderCommand = () => {
  const s = readForm();
  command.textContent = `godi --streams-per-input-device ${s.spid} --verbosity ${s.verbosity} --file-exclude-patterns=${s.fep.join(',')} ${s.mode} ${s.mode !== 'verify' ? `--format ${s.format} ` : ''}${s.sources.join(' ')}${s.mode === 'sealed-copy' ? ` -- ${s.destinations.join(' ')}` : ''}`;
};
form.addEventListener('input', renderCommand);
form.addEventListener('submit', async event => {
  event.preventDefault();
  output.textContent = '';
  const response = await fetch('/api/v1/state',{method:'POST',headers,body:JSON.stringify(readForm())});
  if (!response.ok) output.textContent += `${await response.text()}\n`;
});
document.querySelector('#cancel').onclick = () => fetch('/api/v1/state',{method:'DELETE',headers});

fetch('/api/v1/state',{headers}).then(r => r.json()).then(value => {
  state = value;
  for (const name of ['mode','verbosity','format','spid','spod']) form.elements[name].value = value[name];
  form.elements.fep.value = value.fep.join(',');
  form.elements.sources.value = value.sources.join('\n');
  form.elements.destinations.value = value.destinations.join('\n');
  renderCommand();
  const socket = new WebSocket(`${location.protocol === 'https:' ? 'wss' : 'ws'}://${location.host}${value.socketURL}`);
  socket.onmessage = event => {
    const value = JSON.parse(event.data);
    if (value.state === 1) {
      output.textContent += `${value.message}\n`;
      output.scrollTop = output.scrollHeight;
    }
  };
});
