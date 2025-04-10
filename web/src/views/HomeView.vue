<template>
  <div class="w-full h-full bg-purple-100">
    <loading v-show="!init" />
    <div class="h-full w-full p-4 overflow-scroll" v-show="init">
      <div class="head flex justify-between items-center">
        <div class="title text-3xl font-[400] text-gray-700">{{ data.pos?.adm1 ?? '' }}</div>
        <el-icon size="25px" color="black" @click="handleSearch">
          <img :src="searchSvg">
        </el-icon>
      </div>
      <div class="street mt-10 text-gray-700">
        {{ data.pos?.adm2 ?? '' }}
      </div>
      <div class="wd mt-10 text-gray-700 text-8xl">
        {{ data.weather?.temp ?? '' }}°
      </div>
      <div class="weather text-gray-700 flex gap-2">
        <div class="">
          {{ data.weather?.text ?? '' }}
        </div>
        <i :class="`qi-${data.weather?.icon}`"></i>
        <div>
          最高{{ data.next_weather?.[0].tempMax.padStart(2, ' ') ?? '' }}°
        </div>
        <div>
          最低{{ data.next_weather?.[0].tempMin.padStart(2, ' ') ?? '' }}°
        </div>
      </div>
      <div class="air mt-4 bg-gray-500/70 w-fit p-1 pl-2 pr-2 rounded-full flex items-center">
        <div class="icon w-4 h-4">
          <img :src="leafSvg" alt="leaf">
        </div>
        <div class="text-xs text-white ml-2">
          {{ data.air_quality?.category ?? '' }} {{ data.air_quality?.aqi ?? '' }}
        </div>
      </div>
      <div class="report mt-10 p-5 backdrop-filter bg-white/40 rounded-lg shadow-md">
        <div class="title text-gray-700 font-semibold">
          3日天气预报
        </div>
        <div class="days mt-4 flex flex-col gap-3 text-gray-700">
          <div v-for="i in 3" :key="i" class="day flex text-sm items-center">
            <div class="data">
              {{ getDay(i) }}
            </div>
            <div class="weather ml-4">
              {{ data.next_weather?.[i - 1].textDay ?? '' }}
            </div>
            <div class="flex-1"></div>
            <div class="icon h-6 w-6">
              <i :class="`qi-${data.next_weather?.[i - 1].iconDay}`"></i>
            </div>
            <div class="min">
              {{ data.next_weather?.[i - 1].tempMin.padStart(2, '&nbsp') ?? '' }}°
            </div>
            <div class="line w-20 h-2 ml-2 bg-gradient-to-r from-green-400 via-yellow-500 to-red-400 rounded-full">
            </div>
            <div class="max ml-2">
              {{ data.next_weather?.[i - 1].tempMax.padStart(2, '&nbsp') ?? '' }}°
            </div>
          </div>
        </div>
      </div>
      <div class="air mt-10 p-5 backdrop-filter bg-white/40 rounded-lg shadow-md">
        <div class="title text-gray-700 font-semibold">
          空气质量
        </div>
        <div :class="`row mt-4 ${getColor(Number(data.air_quality?.aqi) ?? 0)}`">
          <span class="text-3xl">{{ data.air_quality?.aqi ?? '' }}</span>
          <span class="font-semibold ml-2">{{ data.air_quality?.category ?? '' }}</span>
        </div>
        <div class="desc text-gray-700 text-xs mt-2">
          {{ data.air_quality?.health.effect ?? '' }}
        </div>
      </div>
      <div class="temperature mt-10 p-5 backdrop-filter bg-white/40 rounded-lg shadow-md">
        <div class="title text-gray-700 font-semibold">
          温度曲线
        </div>
        <div class="w-full overflow-x-auto" ref="tem">
          <div ref="chart" id="chart" class="w-250 h-40"></div>
        </div>
      </div>
      <div class="temperature mt-10 p-5 backdrop-filter bg-white/40 rounded-lg shadow-md">
        <div class="title text-gray-700 font-semibold">
          湿度曲线
        </div>
        <div class="w-full overflow-x-auto" ref="tem2">
          <div ref="chart2" id="chart2" class="w-250 h-40"></div>
        </div>
      </div>
    </div>
    <el-dialog v-model="searchDialog" title="查找城市天气" width="300">
      <el-autocomplete v-model="query" placeholder="请输入城市" clearable :trigger-on-focus="false" :fetch-suggestions="querySearch"
        @select="handleSelect">
      </el-autocomplete>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import leafSvg from '../assets/leaf.svg'
import searchSvg from '../assets/search.svg'
import pinyin from "js-pinyin";
import { onMounted, ref, reactive, watch, nextTick } from 'vue'
import { getGeolocationPromise } from '../ts/api.ts'
import { renderChart, renderChart2 } from '../ts/renderCharts.ts'
import Loading from '../components/LoadingCom.vue';
const query = ref('')
const chart = ref(null)
const chart2 = ref(null)
const tem = ref<HTMLDivElement | null>(null)
const tem2 = ref<HTMLDivElement | null>(null)
const searchDialog = ref(false);
const data = reactive<dataProps>({})
const url = "https://weather.hxzzz.asia/backend/api/weather"
const url2 = "https://weather.hxzzz.asia/backend/api/position"
const init = ref(false);
const isMount = ref(false);
const x = ref<string[]>([])
const y = ref<number[]>([])
const x2 = ref<string[]>([])
const y2 = ref<number[]>([])
const lat = ref('')
const lon = ref('')


const querySearch = (queryString: string, cb: any) => {
  let lists = <{ value: string, lat: string, lon: string }[]>[]
  fetch(url2, {
    method: "POST",
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify({
      "position": pinyin.getFullChars(queryString).toLocaleLowerCase()
    })
  })
    .then(res => res.json())
    .then(res => {
      console.log(res.data)
      res.data.pos_list.forEach((item: { adm1: string, adm2: string, name: string, lat: string, lon: string }) => {
        lists.push({ value: `${item.adm1} ${item.adm2} ${item.name}`, lat: item.lat, lon: item.lon })
      })
      cb(lists)
    })
    .catch(err => {
      console.error(err);
      cb([])
    })
}

const handleSelect = (item: Record<string, any>) => {
  lat.value = item.lat;
  lon.value = item.lon;
  searchDialog.value = false;
}


interface dataProps {
  pos?: {
    adm1: string,
    adm2: string,
  },
  weather?: {
    temp: string,
    text: string,
    icon: string,
  },
  next_weather?: [
    {
      textDay: string,
      tempMin: string,
      tempMax: string,
      iconDay: string,
    }
  ],
  air_quality?: {
    category: string,
    aqi: string,
    health: {
      effect: string,
    }
  }
}


const upDateweather = async () => {
  try {
    const response = await fetch(url, {
      method: 'POST',
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        "longitude": `${lon.value}`,
        "latitude": `${lat.value}`,
      }),
    })
    const newData = await response.json()
    Object.assign(data, newData.data);
    init.value = true;
    x.value = [];
    y.value = [];
    x2.value = [];
    y2.value = [];
    newData.data.hourly_weather.forEach((element: { fxTime: string, temp: string, humidity: string, text: string }) => {
      let date = new Date(element.fxTime);
      x.value.push(`${element.text}\n${String((date.getUTCHours() + 8) % 24).padStart(2, '0')}:00`)
      y.value.push(Number(element.temp))
      x2.value.push(`${String((date.getUTCHours() + 8) % 24).padStart(2, '0')}:00`)
      y2.value.push(Number(element.humidity))
    });
    renderChart(chart, x, y);
    renderChart2(chart2, x2, y2);
  } catch (err) {
    console.error(err);
  }
}
watch(lat, (newLat, _) => {
  console.log("change!!!!!!!!!!!!!")
  if (newLat !== '' && lon.value !== '') {
    upDateweather();
  }
})

const getPos = async() => {
  const position = await getGeolocationPromise();
  lat.value = `${position.latitude}`;
  lon.value = `${position.longitude}`
}

getPos()


function getDay(i: number): string {
  if (i == 1) {
    return "今天"
  } else if (i == 2) {
    return "明天"
  } else {
    return "后天"
  }
}

onMounted(() => {
  isMount.value = true;
})

watch(init, (newData, _) => {
  if (newData == true && isMount.value === true) {
    nextTick(() => {
      renderChart(chart, x, y);
      renderChart2(chart2, x2, y2);

    })
  }
})

watch(isMount, (newData, _) => {
  if (newData === true && init.value === true) {
    nextTick(() => {
      renderChart(chart, x, y);
      renderChart2(chart2, x2, y2);
    })
  }
})

function getColor(i: number): string {
  if (i < 100) {
    return "text-green-700"
  } else if (i < 200) {
    return "text-yellow-600"
  } else {
    return "text-red-700"
  }
}

onMounted(() => {
  const scrollContainer = tem.value;
  const scrollContainer2 = tem2.value;
  if (scrollContainer !== null) {
    scrollContainer.addEventListener('wheel', function (e) {
      e.preventDefault();
      this.scrollLeft += e.deltaY;
    });
  }
  if (scrollContainer2 !== null) {
    scrollContainer2.addEventListener('wheel', function (e) {
      e.preventDefault();
      this.scrollLeft += e.deltaY;
    });
  }
})

const handleSearch = () => {
  searchDialog.value = true;
}
</script>