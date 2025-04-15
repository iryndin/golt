const PossibleStepMillis = [100, 200, 500, 1000, 2000, 5000 ];

function createBarCharData(data) {
  let res;
  for (let i=0; i<PossibleStepMillis.length; i++) {
    const stepMs = PossibleStepMillis[i];
    res = createBarCharDataWithStep(data, stepMs);
    if (res.length < 50) {
      break;
    }
  }
  return res;
}

function createBarCharDataWithStep(data, barChartStepMs) {
  let startMs = 0;
  let result = [];

  let maxMs = 0;
  for (let i=0; i<data.length; i++) {
    if (data[i].elapsedMs > maxMs) {
      maxMs = data[i].elapsedMs;
    }
  }

  while (startMs < maxMs) {
    const endMs = startMs + barChartStepMs;
    let vals = countValuesInInterval(data, startMs, endMs);
    result.push({ts: startMs, n: vals});
    startMs = endMs;
  }

  return result;
}

function countValuesInInterval(data, start, end) {
  let vals = 0;
  for (let i=0; i<data.length; i++) {
    if (start <= data[i].elapsedMs && data[i].elapsedMs < end) {
      vals++;
    }
  }
  return vals;
}
