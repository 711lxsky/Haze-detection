<template>
    <div class="box w-full h-full flex justify-center items-center">
        <div class="w-50 h-50">
            <svg viewBox="0 -960 960 960">
                <g stroke="#06f" fill="none" fill-rule="evenodd" stroke-linecap="round" stroke-linejoin="round"
                    stroke-width="10">
                    <path class="line"
                        d="M260-160q-91 0-155.5-63T40-377q0-78 47-139t123-78q25-92 100-149t170-57q117 0 198.5 81.5T760-520q69 8 114.5 59.5T920-340q0 75-52.5 127.5T740-160H260Zm0-80h480q42 0 71-29t29-71q0-42-29-71t-71-29h-60v-80q0-83-58.5-141.5T480-720q-83 0-141.5 58.5T280-520h-20q-58 0-99 41t-41 99q0 58 41 99t99 41Zm220-240Z" />
                </g>
            </svg>
        </div>
    </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue';
import { animate, svg, stagger } from 'animejs';
onMounted(():void => {
    animate(svg.createDrawable('.line'), {
        draw: ['0 0', '0 1', '1 1'],
        ease: 'inOutQuad',
        duration: 2000,
        delay: stagger(100),
        loop: true
    });
})

</script>