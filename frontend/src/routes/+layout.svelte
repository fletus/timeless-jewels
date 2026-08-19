<script lang="ts">
  import '../app.scss';
  import '../wasm_exec.js';
  import { assets } from '$app/paths';
  import { browser } from '$app/environment';
  import { loadSkillTree } from '../lib/skill_tree';
  import { syncWrap } from '../lib/worker';
  import { initializeCrystalline } from '../lib/types';

  let wasmLoading = true;
  let maxTotalWeight = browser ? localStorage.getItem('maxTotalWeight') ?? '' : '';

  $: if (browser) {
    localStorage.setItem('maxTotalWeight', maxTotalWeight);
  }

  // eslint-disable-next-line no-undef
  const go = new Go();

  if (browser) {
    fetch(assets + '/calculator.wasm')
      .then((data) => data.arrayBuffer())
      .then((data) => {
        WebAssembly.instantiate(data, go.importObject).then((result) => {
          go.run(result.instance);
          wasmLoading = false;
          initializeCrystalline();
          loadSkillTree();
        });

        syncWrap.boot(data);
      });
  }
</script>

{#if wasmLoading}
  <div class="flex flex-row justify-center h-screen">
    <div class="flex flex-col">
      <div class="py-10 flex flex-col justify-between">
        <div>
          <h1 class="text-white mb-10 text-center">Timeless Calculator</h1>

          <h2 class="text-center">Loading...</h2>
        </div>
      </div>
    </div>
  </div>
{:else}
  <slot />
  <div class="fixed bottom-2 left-2 z-50 flex flex-row items-center gap-2 bg-black/80 backdrop-blur-sm themed rounded p-2">
    <div class="min-w-fit">Max Total Weight:</div>
    <input type="number" min="0" bind:value={maxTotalWeight} class="w-24" placeholder="No limit" />
    <span class="text-sm text-neutral-400">blank = no limit</span>
  </div>
{/if}
