package com.greencompass.feature.onboarding

import androidx.compose.foundation.BorderStroke
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.ArrowBack
import androidx.compose.material.icons.filled.Check
import androidx.compose.material.icons.filled.MyLocation
import androidx.compose.material.icons.filled.Search
import androidx.compose.material.icons.outlined.Map
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.unit.dp
import com.greencompass.core.ui.*

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun LocationSelectionScreen(
    onBack: () -> Unit,
    onUseLocation: () -> Unit,
    onSearchPlace: () -> Unit,
    onChooseOnMap: () -> Unit,
    onContinue: () -> Unit
) {
    var selected by remember { mutableStateOf(true) }

    GreenCompassScaffold(
        title = "",
        navigationIcon = {
            IconButton(onClick = onBack) {
                Icon(Icons.AutoMirrored.Filled.ArrowBack, contentDescription = "Back", tint = GreenCompassColors.Charcoal)
            }
        }
    ) { paddingValues ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(paddingValues)
                .padding(horizontal = AppSpacing.lg)
        ) {
            Text(
                text = "Choose a place",
                style = GreenCompassTypography.headlineLarge,
                color = GreenCompassColors.Charcoal,
                modifier = Modifier.padding(bottom = AppSpacing.xs)
            )

            Text(
                text = "See updates for a place that matters to you.",
                style = GreenCompassTypography.bodyMedium,
                color = GreenCompassColors.MutedText,
                modifier = Modifier.padding(bottom = AppSpacing.xl)
            )

            SecondaryButton(
                text = "Use my location",
                onClick = onUseLocation,
                modifier = Modifier.padding(bottom = AppSpacing.sm)
            )

            SecondaryButton(
                text = "Search for a place",
                onClick = onSearchPlace,
                modifier = Modifier.padding(bottom = AppSpacing.sm)
            )

            SecondaryButton(
                text = "Choose on map",
                onClick = onChooseOnMap,
                modifier = Modifier.padding(bottom = AppSpacing.xxl)
            )

            Text(
                text = "Selected place",
                style = GreenCompassTypography.labelMedium,
                color = GreenCompassColors.MutedText,
                modifier = Modifier.padding(bottom = AppSpacing.xs)
            )

            Surface(
                modifier = Modifier.fillMaxWidth().clickable { selected = !selected },
                shape = RoundedCornerShape(12.dp),
                color = if (selected) GreenCompassColors.SoftSage else Color.White,
                border = BorderStroke(1.dp, if (selected) GreenCompassColors.ForestGreen else GreenCompassColors.Stone)
            ) {
                Row(
                    modifier = Modifier.fillMaxWidth().padding(AppSpacing.md),
                    horizontalArrangement = Arrangement.SpaceBetween,
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Text(text = "Lower Valley", style = GreenCompassTypography.bodyLarge, color = GreenCompassColors.Charcoal)
                    if (selected) Icon(Icons.Default.Check, contentDescription = "Selected", tint = GreenCompassColors.ForestGreen)
                }
            }

            Spacer(modifier = Modifier.weight(1f))

            PrimaryButton(
                text = "Continue",
                onClick = onContinue
            )

            Spacer(modifier = Modifier.height(AppSpacing.xxl))
        }
    }
}
