<template>
  <div class="container mx-auto p-6">
    <div v-if="recipe" class="bg-white p-6 rounded shadow">
      <img :src="recipe.images" class="w-full h-64 object-cover rounded mb-4" />
      <h1 class="text-2xl font-bold mb-2">{{ recipe.name }}</h1>
      <p class="text-gray-700 mb-4">{{ recipe.description }}</p>

      <div class="mb-4">
        <strong>烹饪时间：</strong>{{ formatCookTime(recipe.cook_time) }}
      </div>

      <div class="whitespace-pre-line text-gray-800">
        {{ recipe.process }}
      </div>
    </div>
    <div v-else class="text-center text-gray-500 mt-10">正在加载数据...</div>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  data() {
    return {
      recipe: null,
    };
  },
  mounted() {
    this.fetchRecipeDetail();
  },
  methods: {
    async fetchRecipeDetail() {
      const id = this.$route.params.id;
      try {
        const res = await axios.get(`http://127.0.0.1:3000/recipe/detail/${id}`);
        this.recipe = res.data.recipe;
      } catch (error) {
        console.error('获取食谱详情失败:', error);
      }
    },
    formatCookTime(timeStr) {
      const parts = timeStr.split('-');
      return `${parts[1]}分钟`;
    },
  },
};
</script>