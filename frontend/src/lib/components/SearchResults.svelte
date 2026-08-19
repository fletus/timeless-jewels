<script lang="ts">
  import type { SearchResults, SearchWithSeed } from '../skill_tree';
  import SearchResult from './SearchResult.svelte';
  import VirtualList from 'svelte-tiny-virtual-list';

  export let searchResults: SearchResults;
  export let highlight: (newSeed: number, passives: number[]) => void;
  export let groupResults = true;
  export let jewel: number;
  // null means "any conqueror" — see constructQuery in trade.ts
  export let conqueror: string | null;
  export let platform: string;
  export let league: string;
  export let isLegacyTradersMode = false;

  const computeSize = (r: SearchWithSeed) =>
    8 + 48 + r.skills.reduce((o, s) => o + 32 + Object.keys(s.stats).length * 24, 0);

  type SeedRange = {
    min: number;
    max: number;
    step: number;
  };

  // Keep these aligned with data.TimelessJewelSeedRanges.
  // Elegant Hubris only has valid seeds in increments of 20.
  const seedRanges: Record<number, SeedRange> = {
    1: { min: 100, max: 8000, step: 1 },
    2: { min: 10000, max: 18000, step: 1 },
    3: { min: 500, max: 8000, step: 1 },
    4: { min: 2000, max: 10000, step: 1 },
    5: { min: 2000, max: 160000, step: 20 },
    6: { min: 100, max: 8000, step: 1 }
  };

  const buildSeedRange = ({ min, max, step }: SeedRange): number[] => {
    const seeds: number[] = [];
    for (let seed = min; seed <= max; seed += step) {
      seeds.push(seed);
    }
    return seeds;
  };

  let expandedGroup: string | number = '';
  let copiedList: 'matched' | 'inverted' | '' = '';

  $: matchedSeeds = Array.from(new Set(searchResults.raw.map((result) => result.seed))).sort((a, b) => a - b);
  $: matchedSeedSet = new Set(matchedSeeds);
  $: selectedSeedRange = seedRanges[jewel];
  $: invertedSeeds = selectedSeedRange
    ? buildSeedRange(selectedSeedRange).filter((seed) => !matchedSeedSet.has(seed))
    : [];
  $: matchedSeedText = matchedSeeds.join('\n');
  $: invertedSeedText = invertedSeeds.join('\n');

  const copySeeds = async (list: 'matched' | 'inverted', text: string) => {
    await navigator.clipboard.writeText(text);
    copiedList = list;
    setTimeout(() => {
      if (copiedList === list) {
        copiedList = '';
      }
    }, 1500);
  };
</script>

<div class="grid grid-cols-1 xl:grid-cols-2 gap-2 mb-4">
  <div class="flex flex-col bg-neutral-500/20 rounded p-2 min-h-0">
    <div class="flex flex-row items-center justify-between gap-2 mb-2">
      <div class="font-semibold">Matched Seeds [{matchedSeeds.length}]</div>
      <button class="px-3 py-1 bg-neutral-500/40 rounded" on:click={() => copySeeds('matched', matchedSeedText)}>
        {copiedList === 'matched' ? 'Copied' : 'Copy All'}
      </button>
    </div>
    <textarea
      class="w-full h-56 font-mono text-sm resize-y"
      readonly
      value={matchedSeedText}
      aria-label="Matched seed list"></textarea>
  </div>

  <div class="flex flex-col bg-neutral-500/20 rounded p-2 min-h-0">
    <div class="flex flex-row items-center justify-between gap-2 mb-2">
      <div class="font-semibold">Inverted Seeds [{invertedSeeds.length}]</div>
      <button class="px-3 py-1 bg-neutral-500/40 rounded" on:click={() => copySeeds('inverted', invertedSeedText)}>
        {copiedList === 'inverted' ? 'Copied' : 'Copy All'}
      </button>
    </div>
    <textarea
      class="w-full h-56 font-mono text-sm resize-y"
      readonly
      value={invertedSeedText}
      aria-label="Inverted seed list"></textarea>
  </div>
</div>

{#if groupResults}
  <div class="flex flex-col overflow-auto">
    {#each Object.keys(searchResults.grouped)
      .map((x) => parseInt(x))
      .sort((a, b) => a - b)
      .reverse() as k}
      <button
        class="text-lg w-full p-2 px-4 bg-neutral-500/30 rounded flex flex-row justify-between mb-2"
        on:click={() => (expandedGroup = expandedGroup === k ? '' : k)}>
        <span>
          {k} Match{k > 1 ? 'es' : ''} [{searchResults.grouped[k].length}]
        </span>
        <span>
          {expandedGroup === k ? '^' : 'V'}
        </span>
      </button>

      {#if expandedGroup === k}
        <div class="flex flex-col overflow-auto min-h-[200px] mb-2">
          <VirtualList
            height="auto"
            overscanCount={10}
            itemCount={searchResults.grouped[k].length}
            itemSize={searchResults.grouped[k].map(computeSize)}>
            <div slot="item" let:index let:style {style}>
              <SearchResult set={searchResults.grouped[k][index]} {highlight} {jewel} {conqueror} {platform} {league} {isLegacyTradersMode} />
            </div>
          </VirtualList>
        </div>
      {/if}
    {/each}
  </div>
{:else}
  <div class="mt-4 flex flex-col overflow-auto">
    <VirtualList
      height="auto"
      overscanCount={15}
      itemCount={searchResults.raw.length}
      itemSize={searchResults.raw.map(computeSize)}>
      <div slot="item" let:index let:style {style}>
        <SearchResult set={searchResults.raw[index]} {highlight} {jewel} {conqueror} {platform} {league} {isLegacyTradersMode} />
      </div>
    </VirtualList>
  </div>
{/if}
