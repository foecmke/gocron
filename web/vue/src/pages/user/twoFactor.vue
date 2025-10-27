<template>
  <div class="two-factor-container">
    <el-card class="box-card">
      <div slot="header" class="clearfix">
        <span>双因素认证 (2FA)</span>
      </div>
      
      <div v-if="!twoFactorEnabled">
        <el-alert
          title="提示"
          type="info"
          description="启用双因素认证可以大大提高账户安全性。建议所有用户特别是管理员启用此功能。"
          :closable="false"
          show-icon>
        </el-alert>
        
        <el-button 
          type="primary" 
          @click="setup2FA" 
          style="margin-top: 20px;"
          :loading="loading">
          启用2FA
        </el-button>
      </div>

      <div v-else>
        <el-alert
          title="2FA已启用"
          type="success"
          description="您的账户已启用双因素认证保护。"
          :closable="false"
          show-icon>
        </el-alert>
        
        <el-button 
          type="danger" 
          @click="showDisableDialog" 
          style="margin-top: 20px;">
          禁用2FA
        </el-button>
      </div>
    </el-card>

    <el-dialog
      title="启用双因素认证"
      :visible.sync="setupDialogVisible"
      width="500px"
      :close-on-click-modal="false">
      
      <div v-if="qrCode">
        <p>1. 使用认证APP扫描下方二维码：</p>
        <div style="text-align: center; margin: 20px 0;">
          <img :src="qrCode" alt="QR Code" style="width: 200px; height: 200px;">
        </div>
        
        <p>2. 或手动输入密钥：</p>
        <el-input v-model="secret" readonly>
          <el-button slot="append" @click="copySecret">复制</el-button>
        </el-input>
        
        <p style="margin-top: 20px;">3. 输入APP显示的6位验证码：</p>
        <el-input 
          v-model="verifyCode" 
          placeholder="请输入6位验证码"
          maxlength="6"
          @keyup.enter.native="enable2FA">
        </el-input>
      </div>

      <span slot="footer" class="dialog-footer">
        <el-button @click="setupDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="enable2FA" :loading="loading">确定</el-button>
      </span>
    </el-dialog>

    <el-dialog
      title="禁用双因素认证"
      :visible.sync="disableDialogVisible"
      width="400px"
      :close-on-click-modal="false">
      
      <p>请输入认证APP显示的6位验证码以禁用2FA：</p>
      <el-input 
        v-model="disableCode" 
        placeholder="请输入6位验证码"
        maxlength="6"
        @keyup.enter.native="disable2FA">
      </el-input>

      <span slot="footer" class="dialog-footer">
        <el-button @click="disableDialogVisible = false">取消</el-button>
        <el-button type="danger" @click="disable2FA" :loading="loading">确定禁用</el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script>
import userApi from '@/api/user'

export default {
  name: 'TwoFactor',
  data() {
    return {
      twoFactorEnabled: false,
      loading: false,
      setupDialogVisible: false,
      disableDialogVisible: false,
      qrCode: '',
      secret: '',
      verifyCode: '',
      disableCode: ''
    }
  },
  mounted() {
    this.check2FAStatus()
  },
  methods: {
    check2FAStatus () {
      userApi.get2FAStatus((data) => {
        this.twoFactorEnabled = data.enabled
      })
    },

    setup2FA () {
      this.loading = true
      userApi.setup2FA((data) => {
        this.qrCode = data.qr_code
        this.secret = data.secret
        this.setupDialogVisible = true
        this.loading = false
      })
    },

    enable2FA () {
      if (!this.verifyCode || this.verifyCode.length !== 6) {
        this.$message.warning('请输入6位验证码')
        return
      }

      this.loading = true
      userApi.enable2FA(this.secret, this.verifyCode, () => {
        this.$message.success('2FA已启用')
        this.setupDialogVisible = false
        this.twoFactorEnabled = true
        this.verifyCode = ''
        this.loading = false
      })
    },

    showDisableDialog () {
      this.disableCode = ''
      this.disableDialogVisible = true
    },

    disable2FA () {
      if (!this.disableCode || this.disableCode.length !== 6) {
        this.$message.warning('请输入6位验证码')
        return
      }

      this.loading = true
      userApi.disable2FA(this.disableCode, () => {
        this.$message.success('2FA已禁用')
        this.disableDialogVisible = false
        this.twoFactorEnabled = false
        this.disableCode = ''
        this.loading = false
      })
    },
    
    copySecret() {
      const input = document.createElement('input')
      input.value = this.secret
      document.body.appendChild(input)
      input.select()
      document.execCommand('copy')
      document.body.removeChild(input)
      this.$message.success('密钥已复制到剪贴板')
    }
  }
}
</script>

<style scoped>
.two-factor-container {
  padding: 20px;
}

.box-card {
  max-width: 600px;
}
</style>
