package com.greencompass.feature.onboarding

import androidx.compose.foundation.BorderStroke
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.ArrowBack
import androidx.compose.material.icons.filled.Check
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.Icon
import androidx.compose.material3.IconButton
import androidx.compose.material3.Surface
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.unit.dp
import com.greencompass.core.ui.AppSpacing
import com.greencompass.core.ui.GreenCompassColors
import com.greencompass.core.ui.GreenCompassScaffold
import com.greencompass.core.ui.GreenCompassTypography
import com.greencompass.core.ui.PrimaryButton

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun LanguageSelectionScreen(onContinue: () -> Unit, onBack: () -> Unit) {
    var selectedLanguage by remember { mutableStateOf("English") }
    val languages = listOf("English", "Kiswahili", "French", "Arabic")

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
            Text(text = "Choose your language", style = GreenCompassTypography.headlineLarge, color = GreenCompassColors.Charcoal, modifier = Modifier.padding(bottom = AppSpacing.xs))
            Text(text = "You can change this later in Profile.", style = GreenCompassTypography.bodyMedium, color = GreenCompassColors.MutedText, modifier = Modifier.padding(bottom = AppSpacing.xl))

            LazyColumn(verticalArrangement = Arrangement.spacedBy(AppSpacing.sm), modifier = Modifier.weight(1f)) {
                items(languages) { language ->
                    val isSelected = selectedLanguage == language
                    Surface(
                        modifier = Modifier.fillMaxWidth().clickable { selectedLanguage = language },
                        shape = RoundedCornerShape(12.dp),
                        color = if (isSelected) GreenCompassColors.SoftSage else Color.White,
                        border = BorderStroke(1.dp, if (isSelected) GreenCompassColors.ForestGreen else GreenCompassColors.Stone)
                    ) {
                        Row(modifier = Modifier.fillMaxWidth().padding(AppSpacing.md), horizontalArrangement = Arrangement.SpaceBetween, verticalAlignment = Alignment.CenterVertically) {
                            Text(text = language, style = GreenCompassTypography.bodyLarge, color = GreenCompassColors.Charcoal)
                            if (isSelected) Icon(Icons.Default.Check, contentDescription = "Selected", tint = GreenCompassColors.ForestGreen)
                        }
                    }
                }
            }

            Spacer(modifier = Modifier.height(AppSpacing.xl))
            PrimaryButton(text = "Continue", onClick = onContinue)
            Spacer(modifier = Modifier.height(AppSpacing.xxl))
        }
    }
}
