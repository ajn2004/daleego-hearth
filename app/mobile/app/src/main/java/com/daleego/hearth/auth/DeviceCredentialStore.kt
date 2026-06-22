package com.daleego.hearth.auth

import android.content.Context

class DeviceCredentialStore(context: Context) {
    private val prefs = context.getSharedPreferences(PREFS_NAME, Context.MODE_PRIVATE)

    fun getApiKey(): String? {
        return prefs.getString(KEY_API_KEY, null)
    }

    fun saveApiKey(apiKey: String) {
        prefs.edit()
            .putString(KEY_API_KEY, apiKey)
            .apply()
    }

    fun clear() {
        prefs.edit().clear().apply()
    }

    private companion object {
        const val PREFS_NAME = "hearth_device"
        const val KEY_API_KEY = "api_key"
    }
}
