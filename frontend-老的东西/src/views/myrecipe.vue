<template>
  <div class="container mx-auto p-6">
    <h1 class="text-2xl font-bold mb-6">我发布的食谱</h1>

    <!-- 食谱网格 -->
    <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-6">
      <router-link
        v-for="recipe in paginatedRecipes"
        :key="recipe.id"
        :to="`/recipe/detail/${recipe.id}`"
        class="block bg-white rounded shadow hover:shadow-lg transition-shadow"
      >
        <img
          :src="recipe.images"
          alt="Recipe Image"
          class="w-full h-40 object-cover rounded-t"
        />
        <div class="p-4">
          <h2 class="text-lg font-semibold text-gray-800">{{ recipe.name }}</h2>
          <p class="text-sm text-gray-500 mt-1 line-clamp-2">
            {{ recipe.description }}
          </p>
          <div class="mt-2 text-xs text-gray-400 flex justify-between">
            <span>👍 {{ recipe.likes }}</span>
            <span>烹饪时间：{{ formatCookTime(recipe.cook_time) }}</span>
          </div>
        </div>
      </router-link>
    </div>

    <!-- 分页控件 -->
    <div class="mt-8 flex justify-center items-center space-x-4">
      <button
        @click="prevPage"
        :disabled="currentPage === 1"
        class="px-4 py-2 bg-gray-700 text-white rounded disabled:opacity-50"
      >
        上一页
      </button>
      <span class="text-gray-600">第 {{ currentPage }} / {{ totalPages }} 页</span>
      <button
        @click="nextPage"
        :disabled="currentPage === totalPages"
        class="px-4 py-2 bg-gray-700 text-white rounded disabled:opacity-50"
      >
        下一页
      </button>
    </div>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  data() {
    return {
      recipes: [],
      currentPage: 1,
      pageSize: 6,
    };
  },
  mounted() {
    this.fetchUserRecipes();
  },
  methods: {
    async fetchUserRecipes() {
      try {
        const res = await axios.get('http://127.0.0.1:3000/recipe/get');
        this.recipes = res.data.recipes;
      } catch (error) {
        console.error('获取用户食谱失败:', error);
      }
    },
    nextPage() {
      if (this.currentPage < this.totalPages) {
        this.currentPage++;
      }
    },
    prevPage() {
      if (this.currentPage > 1) {
        this.currentPage--;
      }
    },
    formatCookTime(timeStr) {
      // cook_time 格式是 "hh-mm-ss"，只展示 mm 部分即可
      const parts = timeStr.split('-');
      return `${parts[1]}分钟`;
    },
  },
  computed: {
    totalPages() {
      return Math.ceil(this.recipes.length / this.pageSize);
    },
    paginatedRecipes() {
      const start = (this.currentPage - 1) * this.pageSize;
      const end = start + this.pageSize;
      return this.recipes.slice(start, end);
    },
  },
};
</script>