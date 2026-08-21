package com.greencompass.feature.today

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.greencompass.data.repository.GreenCompassRepository
import com.greencompass.domain.model.TodayData
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.update
import kotlinx.coroutines.launch
import javax.inject.Inject

data class TodayUiState(
    val isLoading: Boolean = true,
    val data: TodayData? = null,
    val isOffline: Boolean = false,
    val errorMessage: String? = null
)

@HiltViewModel
class TodayViewModel @Inject constructor(
    private val repository: GreenCompassRepository
) : ViewModel() {

    // Hardcoded to "1" (Lower Valley) for demo. 
    // In a real app, this would come from a PlaceManager/Settings repository.
    private val currentPlaceId = "1" 

    private val _uiState = MutableStateFlow(TodayUiState())
    val uiState: StateFlow<TodayUiState> = _uiState.asStateFlow()

    init {
        observeData()
    }

    private fun observeData() {
        viewModelScope.launch {
            repository.observeToday(currentPlaceId).collect { data ->
                _uiState.update {
                    it.copy(
                        isLoading = false,
                        data = data
                    )
                }
            }
        }
    }

    fun toggleOfflineMode() {
        _uiState.update { it.copy(isOffline = !it.isOffline) }
    }
    
    fun retry() {
        _uiState.update { it.copy(isLoading = true, errorMessage = null) }
        observeData()
    }
}
