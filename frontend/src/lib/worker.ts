import { browser } from '$app/environment';
import SyncWorker from './sync_worker?worker';
import * as Comlink from 'comlink';
import type { WorkerType } from './sync_worker';

function getWorker() {
  console.log('Creating sync worker');
  const theWorker = new SyncWorker();
  const remote = Comlink.wrap<WorkerType>(theWorker);
  const obj = {
    boot: remote.boot,
    search: (args: Parameters<WorkerType['search']>[0], callback: Parameters<WorkerType['search']>[1]) => {
      const storedMax = localStorage.getItem('maxTotalWeight');
      return remote.search(
        {
          ...args,
          maxTotalWeight: storedMax === null || storedMax === '' ? -1 : Number(storedMax)
        },
        callback
      );
    }
  } as typeof remote;
  return { syncWorker: theWorker, syncWrap: obj };
}

export const { syncWorker, syncWrap } = browser ? getWorker() : { syncWorker: null, syncWrap: null };
