package com.greencompass.feature.places

import androidx.compose.foundation.layout.*
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.ArrowBack
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import com.greencompass.core.ui.*

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun AddPlaceScreen(
    onBack: () -> Unit,
    onUseLocation: () -> Unit,
    onSearch: () -> Unit,
    onChooseOnMap: () -> Unit
) {
    GreenCompassScaffold(
        title = "",
        navigationIcon = {
            IconButton(onClick = onBack) {
                Icon(Icons.AutoMirrored.Filled.ArrowBack, contentDescription = "Back", tint = GreenCompassColors.Charcoal)
            }
        }
    ) { paddingValues ->
        Column(
            modifier = Modifier.fillMaxSize().padding(paddingValues).padding(horizontal = AppSpacing.lg)
        ) {
            Spacer(modifier = Modifier.height(AppSpacing.xxl))
            Text(text = "Add a place", style = GreenCompassTypography.headlineLarge, color = GreenCompassColors.Charcoal, modifier = Modifier.padding(bottom = AppSpacing.xl))

            SecondaryButton(text = "Use my location", onClick = onUseLocation, modifier = Modifier.padding(bottom = AppSpacing.sm))
            SecondaryButton(text = "Search for a place", onClick = onSearch, modifier = Modifier.padding(bottom = AppSpacing.sm))
            SecondaryButton(text = "Choose on map", onClick = onChooseOnMap, modifier = Modifier.padding(bottom = AppSpacing.xxl))

            Spacer(modifier = Modifier.weight(1f))
        }
    }
}
