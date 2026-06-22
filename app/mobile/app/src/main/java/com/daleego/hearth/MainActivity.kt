package com.daleego.hearth

import android.os.Build
import android.os.Bundle
import android.Manifest
import android.content.pm.PackageManager
import androidx.activity.result.contract.ActivityResultContracts
import androidx.core.content.ContextCompat
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
import com.daleego.hearth.location.LocationProvider
import com.daleego.hearth.api.PairingRequest
import com.daleego.hearth.api.readableApiError
import com.daleego.hearth.auth.DeviceCredentialStore
import com.daleego.hearth.ui.HomeScreen
import com.daleego.hearth.ui.PairingScreen
import com.daleego.hearth.BuildConfig
import kotlinx.coroutines.launch
import java.time.Instant

private const val BASE_URL = BuildConfig.HEARTH_BASE_URL

class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)

        val credentials = DeviceCredentialStore(this)
        val api: HearthApi = HearthApiClient.create(
            baseUrl = BASE_URL,
            getDeviceApiKey = {credentials.getApiKey()}
            )

        val locationProvider = LocationProvider(this)

        val locationPermissionLauncher = registerForActivityResult(
            ActivityResultContracts.RequestMultiplePermissions()
        ) { permissions ->
            val fineGranted = permissions[Manifest.permission.ACCESS_FINE_LOCATION] == true
            val coarseGranted = permissions[Manifest.permission.ACCESS_COARSE_LOCATION] == true

            if (!fineGranted && !coarseGranted) {
                // You can reflect this in UI later
            }
        }

        fun hasLocationPermission(): Boolean {
            val fineGranted = ContextCompat.checkSelfPermission(
                this,
                Manifest.permission.ACCESS_FINE_LOCATION
            ) == PackageManager.PERMISSION_GRANTED

            val coarseGranted = ContextCompat.checkSelfPermission(
                this,
                Manifest.permission.ACCESS_COARSE_LOCATION
            ) == PackageManager.PERMISSION_GRANTED

            return fineGranted || coarseGranted
        }
        setContent {
            val scope = rememberCoroutineScope()
            var apiKey by remember { mutableStateOf(credentials.getApiKey()) }
            var deviceName by remember { mutableStateOf<String?>(null) }
            var isPairing by remember { mutableStateOf(false) }
            var pairingError by remember { mutableStateOf<String?>(null) }
            var isCheckInLoading by remember { mutableStateOf(false) }
            var checkInStatus by remember { mutableStateOf<String?>(null) }
            var isTracking by remember { mutableStateOf(false) }

            suspend fun sendCurrentLocation() {
                if (!hasLocationPermission()) {
                    locationPermissionLauncher.launch(
                        arrayOf(
                            android.Manifest.permission.ACCESS_FINE_LOCATION,
                            android.Manifest.permission.ACCESS_COARSE_LOCATION
                        )
                    )
                    throw IllegalStateException("Location permission required.")
                }

                val location = locationProvider.getCurrentLocation()
                    ?: throw IllegalStateException("No location available yet.")

                api.reportLocation(
                    LocationRequest(
                        lat = location.latitude,
                        lng = location.longitude,
                        accuracy = location.accuracy.toDouble(),
                        recordedAt = java.time.Instant.ofEpochMilli(location.time).toString()
                    )
                )
            }

            androidx.compose.runtime.LaunchedEffect(apiKey, isTracking) {
                if (apiKey == null || !isTracking) return@LaunchedEffect

                while (isTracking) {
                    try {
                        sendCurrentLocation()
                        checkInStatus = "Location sent at ${java.time.Instant.now()}"
                    } catch (throwable: Exception) {
                        checkInStatus = readableApiError(throwable)
                    }

                    kotlinx.coroutines.delay(10 * 1000L) // 5 minutes
                }
            }

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
                    isTracking = isTracking,
                    onStartTracking = {
                        isTracking = true
                    },
                    onStopTracking = {
                        isTracking = false
                    },
                    onSendTestCheckIn = {
                        scope.launch {
                            isCheckInLoading = true
                            checkInStatus = null

                            try {
                                if (!hasLocationPermission()) {
                                    locationPermissionLauncher.launch(
                                        arrayOf(
                                            Manifest.permission.ACCESS_FINE_LOCATION,
                                            Manifest.permission.ACCESS_COARSE_LOCATION
                                        )
                                    )

                                    checkInStatus = "Location permission requested."
                                    return@launch
                                }

                                val location = locationProvider.getCurrentLocation()

                                if (location == null) {
                                    checkInStatus = "No location available yet. Try again in a moment."
                                    return@launch
                                }

                                api.reportLocation(
                                    LocationRequest(
                                        lat = location.latitude,
                                        lng = location.longitude,
                                        accuracy = location.accuracy.toDouble(),
                                        recordedAt = Instant.ofEpochMilli(location.time).toString()
                                    )
                                )

                                checkInStatus = "Location sent."
                            } catch (throwable: Exception) {
                                checkInStatus = readableApiError(throwable)
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
