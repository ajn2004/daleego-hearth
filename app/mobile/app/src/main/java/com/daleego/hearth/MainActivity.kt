package com.daleego.hearth

import android.os.Build
import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.rememberCoroutineScope
import androidx.compose.runtime.setValue
import com.daleego.hearth.api.HearthApi
import com.daleego.hearth.api.HearthApiClient
import com.daleego.hearth.api.LocationRequest
import com.daleego.hearth.api.PairingRequest
import com.daleego.hearth.api.readableApiError
import com.daleego.hearth.auth.DeviceCredentialStore
import com.daleego.hearth.ui.HomeScreen
import com.daleego.hearth.ui.PairingScreen
import kotlinx.coroutines.launch

private const val BASE_URL = "http://192.168.0.252:8080/"

class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)

        val credentials = DeviceCredentialStore(this)
        val api: HearthApi = HearthApiClient.create(BASE_URL)

        setContent {
            val scope = rememberCoroutineScope()
            var apiKey by remember { mutableStateOf(credentials.getApiKey()) }
            var deviceName by remember { mutableStateOf<String?>(null) }
            var isPairing by remember { mutableStateOf(false) }
            var pairingError by remember { mutableStateOf<String?>(null) }
            var isCheckInLoading by remember { mutableStateOf(false) }
            var checkInStatus by remember { mutableStateOf<String?>(null) }

            if (apiKey == null) {
                PairingScreen(
                    isLoading = isPairing,
                    errorMessage = pairingError,
                    onPair = { pairCode ->
                        scope.launch {
                            isPairing = true
                            pairingError = null

                            try {
                                val response = api.pairDevice(
                                    PairingRequest(
                                        pairCode = pairCode,
                                        platform = "android",
                                        model = Build.MODEL,
                                        modelType = Build.DEVICE
                                    )
                                )

                                credentials.saveApiKey(response.apiKey)
                                apiKey = response.apiKey
                                deviceName = response.device.name
                            } catch (throwable: Exception) {
                                pairingError = readableApiError(throwable)
                            } finally {
                                isPairing = false
                            }
                        }
                    }
                )
            } else {
                HomeScreen(
                    deviceName = deviceName,
                    checkInStatus = checkInStatus,
                    isCheckInLoading = isCheckInLoading,
                    onSendTestCheckIn = {
                        scope.launch {
                            isCheckInLoading = true
                            checkInStatus = null

                            try {
                                api.reportLocation(
                                    LocationRequest(
                                        lat = 0.0,
                                        lng = 0.0,
                                        acc = 0.0
                                    )
                                )
                                checkInStatus = "Test check-in sent."
                            } catch (throwable: Exception) {
                                checkInStatus = throwable.message ?: "Failed to send test check-in."
                            } finally {
                                isCheckInLoading = false
                            }
                        }
                    },
                    onUnpair = {
                        credentials.clear()
                        apiKey = null
                        deviceName = null
                        pairingError = null
                        isCheckInLoading = false
                        checkInStatus = null
                    }
                )
            }
        }
    }
}
