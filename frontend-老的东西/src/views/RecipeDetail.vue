<template>
  <div class="container mx-auto p-6">
    <div v-if="recipe" class="recipecard">
      <h1>{{ recipe.Title }}</h1>
      <div style="width: 80%;border:2px solid #999999;height: 0px; margin:5px auto 0 auto"></div>
      
      <div class="imageStyle">
        <img :src="recipe.images" class="recipe-image"/>
      </div>

      <!-- 重新设计的分享组件 -->
      <div class="share-container">
        <el-button class="share-left"><el-icon><StarFilled/></el-icon> 收藏菜谱</el-button>
        <el-button class="share-left"> <el-icon><Plus /></el-icon>加入专辑</el-button>

        <div class="share-text">分享</div>
        <div class="share-buttons">
          <div class="share-item" @click="shareToQQ" title="分享到QQ空间">
            <img class="share-img" src="../assets/share/share_qq_active.png" alt="QQ">
            <span class="share-label">QQ空间</span>
          </div>
          <div class="share-item" @click="shareToWx" title="分享到微信">
            <img class="share-img" src="../assets/share/share_wx_active.png" alt="微信">
            <span class="share-label">微信</span>
          </div>
          <div class="share-item" @click="shareToWb" title="分享到微博">
            <img class="share-img" src="../assets/share/share_weibo_active.png" alt="微博">
            <span class="share-label">微博</span>
          </div>
          <div class="share-item" @click="shareToMobile" title="手机扫码查看">
            <img class="share-img" src="../assets/share/share_mobile_active.png" alt="手机">
            <span class="share-label">扫码</span>
          </div>
        </div>
      </div>
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